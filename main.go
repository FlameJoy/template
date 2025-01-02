package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"temp/initializers"
	"time"
)

var (
	debug = flag.Bool("debug", false, "debuging code")
)

func main() {
	// Логгирование
	// Создаем логгер с уровнем INFO
	logger := NewLogger(INFO, log.New(os.Stdout, "", log.LstdFlags), false)

	if *debug {
		logger.SetLevel(DEBUG)
	}

	// Port
	initializers.LoadEnv(".env")
	port := os.Getenv("DB_PORT")

	// Создание хендлера
	payloadSignatureKey := os.Getenv("SIGNATURE_KEY")
	proofLifeTimeSec, err := strconv.ParseInt(os.Getenv("PAYLOAD_LIFETIME_SEC"), 10, 64)
	if err != nil {
		logger.Fatal("error: %s", err)
	}

	h := newHandler(payloadSignatureKey, time.Duration(proofLifeTimeSec)*time.Second, logger)

	// Регистрация маршрутов
	mux := http.NewServeMux()
	registerHandlers(mux, h, logger)

	// Настройка сервера
	server := http.Server{
		Addr:         ":" + port,
		Handler:      recoveryMiddleware(logger)(loggingMiddleware(logger)(mux)),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Завершение работы сервера по сигналу
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	go func() {
		logger.Info("Server is running at http://0.0.0.0%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Could not listen on :8080: %v\n", err)
		}
	}()

	<-done
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed: %v\n", err)
	}
	logger.Info("Server exited properly")
}
