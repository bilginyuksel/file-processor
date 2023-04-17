package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const shutdownTimeoutDuration = time.Second * 5

func main() {
	zap.ReplaceGlobals(zap.NewExample())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	e := echo.New()
	fileprocrRestHandler := &fileprocr.RestHandler{}
	fileprocrRestHandler.RegisterRoutes(e)

	zap.L().Info("Starting echo server")

	go func() {
		if err := e.Start(":8010"); err != nil {
			quit <- os.Interrupt
		}
	}()

	<-quit

	zap.L().Info("Shutting down echo server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutDuration)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		zap.L().Error("Error while shutting down echo server", zap.Error(err))
	}
}
