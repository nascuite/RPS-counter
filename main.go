package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// hits stores the number of HTTP requests received during the current 1‑second window.
var hits atomic.Uint64

func main() {
	newLogger, _ := zap.NewProduction()
	defer newLogger.Sync() // flushes buffer, if any
	logger := newLogger.Sugar()

	app := fiber.New() // create a new Fiber instance

	app.Get("/", func(c *fiber.Ctx) error {
		hits.Add(1)
		return c.SendString("ok")
	})

	// pprof handlers (e.g. /debug/pprof/profile) are already on DefaultServeMux.
	app.Get("/debug/pprof/*", adaptor.HTTPHandler(http.DefaultServeMux))

	// Log RPS every second.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			v := hits.Swap(0)
			logger.Infof("RPS: %d\n", v)
		}
	}()

	// Context that is cancelled on SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initiate graceful shutdown once a signal is received.
	go func() {
		<-ctx.Done()
		logger.Infof("shutting down…")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := app.Shutdown(); err != nil {
			logger.Errorf("Ошибка при остановке сервера: %v\n", err)
		}

		<-shutdownCtx.Done()
		logger.Infof("Сервер успешно остановлен.")
	}()

	logger.Infof("listening on :8080")
	logger.Fatal(app.Listen(fmt.Sprintf(":%d", 8080)))
}
