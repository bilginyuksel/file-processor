package main

import (
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/bilginyuksel/file-processor/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const shutdownTimeoutDuration = time.Second * 5

func main() {
	zap.ReplaceGlobals(zap.NewExample())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	e := echo.New()
	fileprocrRestHandler := &fileprocr.RestHandler{}
	fileprocrRestHandler.RegisterRoutes(e)

	go func() {
		zap.L().Info("Starting echo server")
		if err := e.Start(":8010"); err != nil {
			quit <- os.Interrupt
		}
	}()

	grpcSrv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	fileProcrGrpcServer := &fileprocr.GrpcServer{}
	pb.RegisterProcrServer(grpcSrv, fileProcrGrpcServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		zap.L().Info("Starting gRPC server")
		if err := grpcSrv.Serve(lis); err != nil {
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

	zap.L().Info("Shutting down gRPC server")

	grpcSrv.GracefulStop()
}
