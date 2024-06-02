package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"

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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func init() {
	debugMode := common.GetBool("DEBUG_MODE", false)
	if debugMode {
		log.Println("Service RUN on DEBUG mode")
	}

	prometheus.Register(common.GetTotalRequestPrometheus())
	prometheus.Register(common.GetTotalSuccessfulRequestPrometheus())
	prometheus.Register(common.GetTotalFailedRequestPrometheus())
}

func main() {
	logger := logrus.StandardLogger()
	zipkinLogger := log.New(os.Stderr, "twitter-clone", log.Ldate|log.Ltime|log.Llongfile)

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

	tracerProvider, err := initTracer(common.GetString("ZIPKIN_URL", ""), zipkinLogger)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	userMiddleWare := userMiddleware.InitMiddleware()
	e.Use(userMiddleWare.ExtraHeader)

	trace := tracerProvider.Tracer("twitter-clone")

	userRepository := userDBRepository.NewUserRepository(sqlConn, logger, trace)
	userSessionRepository := userSessionRedisRepository.NewUserSessionRepository(redisConn, logger, trace)

	userUsecase := userUsecase.NewUserUsecase(userRepository, userSessionRepository, logger, trace)
	userHandler.NewUserHandler(e, userUsecase, logger, trace)

	tweetRepository := tweetDBRepository.NewTweetRepository(sqlConn, sqlTxConn, logger, trace)
	tweetUsecase := tweetUsecase.NewTweetUsecase(tweetRepository, userSessionRepository, logger, trace)
	tweetHandler.NewTweetHandler(e, tweetUsecase, logger, trace)

	e.GET("/prometheus/metrics", echo.WrapHandler(promhttp.Handler()))

	serverListener := http.Server{
		Addr:    ":9090",
		Handler: e,
	}

	log.Printf("service is listening at port %s", serverListener.Addr)

	twitterCloneServerCrt := common.GetString("TWITTER_CLONE_SERVER_CRT", "")
	twitterCloneServerPK := common.GetString("TWITTER_CLONE_SERVER_PK", "")
	if err := serverListener.ListenAndServeTLS(twitterCloneServerCrt, twitterCloneServerPK); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func createSQLConnectionInstance() (*sql.DB, error) {
	dbHost := common.GetString("DB_HOST", "")
	dbPort := common.GetInt32("DB_HOST_PORT", 0)
	dbUser := common.GetString("DB_USER", "")
	dbPass := common.GetString("DB_PASS", "")
	dbName := common.GetString("DB_NAME", "")

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open(`postgres`, psqlInfo)

	return dbConn, err
}

func createRedisConnectionInstance() *redis.Client {
	redisHost := common.GetString("REDIS_HOST", "")
	redisPort := common.GetInt32("REDIS_HOST_PORT", 0)
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

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(url string, zipkinLogger *log.Logger) (*sdktrace.TracerProvider, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(zipkinLogger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("twitter-clone"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}
