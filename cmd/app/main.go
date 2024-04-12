package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/broker"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/handler"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/repository"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/server"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/service"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/hasher"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/postgres"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/redisdb"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/tokens"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/transactions"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠻⢿⣿⣿⣿⣿⣿⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠇⠄⠸⣿⣿⣿⣿⣿⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠄⠄⢰⣿⣿⣿⣿⣿⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠁⠄⠄⣼⣿⣿⣿⣿⣿⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠁⠄⠄⠄⢸⣿⣿⣿⣿⣿⣿⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠉⣿
// ⣿⠄⠄⠄⠄⠄⣿⠛⠋⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣠⣿
// ⣿⠄⠄⠄⠄⠄⣿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⣿
// ⣿⠄⠄⠄⠄⠄⣿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣴⣿
// ⣿⠄⠄⠄⠄⠄⣿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢈⣿
// ⣿⠄⠄⣴⣶⡄⣿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠰⣿⣿
// ⣿⣤⣤⣭⣯⣤⣿⣿⣿⣷⣶⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣿⣿
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿

func main() {
	if err := loadConfigs(); err != nil {
		logrus.Fatalf("error loading configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env file: %s", err)
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error connect to postgres: %s", err)
	}
	rdb := redisdb.NewRedisDB(redisdb.Config{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		DB:       viper.GetInt("redis.db"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	transactionManager := transactions.NewManager(db)
	ctxTrManager := transactions.NewCtxManager()
	ctxTrGetter := transactions.NewCtxGetter(ctxTrManager)

	repos := repository.NewRepository(repository.Deps{
		DB:        db,
		CtxGetter: ctxTrGetter,
	})

	accessTTL, err := time.ParseDuration(viper.GetString("tokens.accessTTL"))
	if err != nil {
		logrus.Fatalf("invalid accessTTL: %s", err)
	}
	emailTTL, err := time.ParseDuration(viper.GetString("tokens.emailTTL"))
	if err != nil {
		logrus.Fatalf("invalid emailTTL: %s", err)
	}
	tokenManager := tokens.NewTokenManager(tokens.Config{
		SecretKey: os.Getenv("SECRET_KEY"),
		AccessTTL: accessTTL,
		EmailTTL:  emailTTL,
	})

	hasher := hasher.NewHasher(os.Getenv("SALT"))

	broker := broker.NewBroker(broker.Deps{
		RDB: rdb,
	})

	services := service.NewService(service.Deps{
		Repos:              repos,
		TokenManager:       tokenManager,
		RDB:                rdb,
		Hasher:             hasher,
		TransactionManager: transactionManager,
		Broker:             broker,
	})

	handlerDeps := handler.Deps{
		TokenManager: tokenManager,
		Services:     services,
	}
	handler := handler.NewHandler(handlerDeps)
	serverConfig := server.ServerConfig{
		Host: viper.GetString("app.host"),
		Port: viper.GetString("app.port"),
	}
	srv := server.NewServer(serverConfig, handler.InitRoutes())

	go func() {
		if err := srv.Run(); err != nil {
			if !errors.Is(http.ErrServerClosed, err) {
				logrus.Fatalf("error running server: %s", err)
			}
		}
	}()
	logrus.Printf("Server starting...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Server shutting down...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error shutting down server: %s", err)
	}

	logrus.Printf("Server stoped")
}

func loadConfigs() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
