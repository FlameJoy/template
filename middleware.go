package main

import (
	"net/http"
	"time"
)

// Middleware — это функция-обработчик для добавления промежуточной логики
type Middleware func(http.Handler) http.Handler

// Промежуточный слой для логгирования
func loggingMiddleware(logger *CustomLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		})
	}
}

// Промежуточный слой для восстановления из паники
func recoveryMiddleware(logger *CustomLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Warn("Recovered from panic: %v", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Пустышка
func emptyMiddleware(logger *CustomLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("emptyMiddleware is runnning")
			time.Sleep(time.Second)
			logger.Info("emptyMiddleware ended")
			next.ServeHTTP(w, r)
		})
	}
}
