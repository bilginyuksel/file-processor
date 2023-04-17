package main

import (
	"context"
	"fmt"

	"github.com/bilginyuksel/file-processor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewProcrClient(conn)
	res, err := client.Upload(context.Background(), &pb.UploadRequest{})
	fmt.Println(res, err)
}
