package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	userHandler "github.com/alifahsanilsatria/twitter-clone/user/delivery/http"
	userMiddleware "github.com/alifahsanilsatria/twitter-clone/user/delivery/http/middleware"
	userDBRepository "github.com/alifahsanilsatria/twitter-clone/user/repository/db"
	userUsecase "github.com/alifahsanilsatria/twitter-clone/user/usecase"
	userSessionRedisRepository "github.com/alifahsanilsatria/twitter-clone/user_session/repository/redis"
	"github.com/sirupsen/logrus"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
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

	redisConn := createRedisConnectionInstance()

	e := echo.New()

	userMiddleWare := userMiddleware.InitMiddleware()
	e.Use(userMiddleWare.CORS)

	userRepository := userDBRepository.NewUserRepository(sqlConn, logger)
	userSessionRepository := userSessionRedisRepository.NewUserSessionRepository(redisConn, logger)

	userUsecase := userUsecase.NewUserUsecase(userRepository, userSessionRepository, logger)
	userHandler.NewUserHandler(e, userUsecase, logger)

	serverListener := http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: e,
	}

	if err := serverListener.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return
}

func createSQLConnectionInstance() (*sql.DB, error) {
	dbHost := viper.GetString(`database.sql.host`)
	dbPort := viper.GetString(`database.sql.port`)
	dbUser := viper.GetString(`database.sql.user`)
	dbPass := viper.GetString(`database.sql.pass`)
	dbName := viper.GetString(`database.sql.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`postgresql`, dsn)

	return dbConn, err
}

func createRedisConnectionInstance() *redis.Client {
	redisHost := viper.GetString(`database.redis.host`)
	redisPort := viper.GetString(`database.redis.port`)
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	redisPassword := viper.GetString(`database.redis.password`)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return client
}
