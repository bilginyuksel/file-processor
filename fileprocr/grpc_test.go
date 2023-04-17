package fileprocr_test

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/bilginyuksel/file-processor/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestFileprocrGrpcServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping grpc functional test")
	}

	testGrpcServer, url := startTestGrpcServer(t)
	defer testGrpcServer.GracefulStop()

	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	content := []byte("hello world!")
	client := pb.NewProcrClient(conn)
	res, err := client.Upload(context.Background(), &pb.UploadRequest{Data: content})

	t.Log(res)

	assert.NoError(t, err)
	assert.NoError(t, os.Remove(res.GetFilename()))
}

func startTestGrpcServer(t *testing.T) (*grpc.Server, string) {
	lfs := fileprocr.NewLocalFileStorage("")
	require.NoError(t, lfs.Configure())
	svc := fileprocr.NewProcr(1024, lfs)

	procrGrpcServer := fileprocr.NewGrpcServer(svc)
	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterProcrServer(server, procrGrpcServer)

	lis, err := net.Listen("tcp", "")
	require.NoError(t, err)

	go server.Serve(lis)
	return server, lis.Addr().String()
}
