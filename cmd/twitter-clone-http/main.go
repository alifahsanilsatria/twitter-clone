package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/common"
	tweetHandler "github.com/alifahsanilsatria/twitter-clone/tweet/delivery/http"
	tweetDBRepository "github.com/alifahsanilsatria/twitter-clone/tweet/repository/db"
	tweetUsecase "github.com/alifahsanilsatria/twitter-clone/tweet/usecase"
	userHandler "github.com/alifahsanilsatria/twitter-clone/user/delivery/http"
	userMiddleware "github.com/alifahsanilsatria/twitter-clone/user/delivery/http/middleware"
	userDBRepository "github.com/alifahsanilsatria/twitter-clone/user/repository/db"
	userUsecase "github.com/alifahsanilsatria/twitter-clone/user/usecase"
	userSessionRedisRepository "github.com/alifahsanilsatria/twitter-clone/user_session/repository/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	debugMode := common.GetBool("DEBUG_MODE", false)
	if debugMode {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	logger := logrus.StandardLogger()

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	sqlConn, err := createSQLConnectionInstance()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := sqlConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sqlTxConn, err := sqlConn.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	redisConn := createRedisConnectionInstance()

	e := echo.New()

	userMiddleWare := userMiddleware.InitMiddleware()
	e.Use(userMiddleWare.RequestId)

	userRepository := userDBRepository.NewUserRepository(sqlConn, logger)
	userSessionRepository := userSessionRedisRepository.NewUserSessionRepository(redisConn, logger)

	userUsecase := userUsecase.NewUserUsecase(userRepository, userSessionRepository, logger)
	userHandler.NewUserHandler(e, userUsecase, logger)

	tweetRepository := tweetDBRepository.NewTweetRepository(sqlConn, sqlTxConn, logger)
	tweetUsecase := tweetUsecase.NewTweetUsecase(tweetRepository, userSessionRepository, logger)
	tweetHandler.NewTweetHandler(e, tweetUsecase, logger)

	serviceAddress := common.GetString("TWITTER_CLONE_ADDRESS", "")

	serverListener := http.Server{
		Addr:    serviceAddress,
		Handler: e,
	}

	log.Printf("service is listening at port %s", serverListener.Addr)

	if err := serverListener.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func createSQLConnectionInstance() (*sql.DB, error) {
	dbHost := common.GetString("DB_HOST", "")
	dbPort := common.GetInt32("DB_PORT", 0)
	dbUser := common.GetString("DB_USER", "")
	dbPass := common.GetString("DB_PASS", "")
	dbName := common.GetString("DB_NAME", "")

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open(`postgres`, psqlInfo)

	return dbConn, err
}

func createRedisConnectionInstance() *redis.Client {
	redisHost := common.GetString("REDIS_HOST", "")
	redisPort := common.GetInt32("REDIS_PORT", 0)
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisPassword := common.GetString("REDIS_PASS", "")
	redisDBNumber := common.GetInt32("REDIS_DB_NUMBER", 0)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       int(redisDBNumber),
	})

	return client
}
