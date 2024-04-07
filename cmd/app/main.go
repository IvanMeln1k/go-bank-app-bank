package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/handler"
	"github.com/IvanMeln1k/go-bank-app-bank/internal/server"
	"github.com/sirupsen/logrus"
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
	handlerDeps := handler.Deps{}
	handler := handler.NewHandler(handlerDeps)
	serverConfig := server.ServerConfig{
		Host: "localhost",
		Port: "8000",
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
