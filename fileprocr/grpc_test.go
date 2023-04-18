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

func TestFileProcrGrpcServerUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping grpc server tests")
	}

	testCases := []struct {
		desc                        string
		content                     []byte
		extraAssertionsForAsyncTask func(string, *fileprocr.Procr)
	}{
		{
			content: []byte("hello world!"),
			desc:    "File content is not a json, should only save the raw file",
			extraAssertionsForAsyncTask: func(s string, p *fileprocr.Procr) {
				assert.Error(t, <-p.ProcrResultQueue)
				assert.Error(t, os.Remove(s+".json"))
			},
		},
		{
			content: []byte(`{"s": 10, "a": 5}`),
			desc:    "File content is a JSON, should save both raw and json files",
			extraAssertionsForAsyncTask: func(s string, p *fileprocr.Procr) {
				assert.NoError(t, <-p.ProcrResultQueue)
				assert.NoError(t, os.Remove(s+".json"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testGrpcServer, procr, url := startTestGrpcServer(t)
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

			assert.Error(t, <-procr.ProcrResultQueue)
		})
	}
}

func startTestGrpcServer(t *testing.T) (*grpc.Server, *fileprocr.Procr, string) {
	lfs := fileprocr.NewLocalFileStorage("")
	idgen := fileprocr.NewIDGenerator()
	require.NoError(t, lfs.Configure())
	svc := fileprocr.NewProcr(1024, lfs, idgen)

	procrGrpcServer := fileprocr.NewGrpcServer(svc)
	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterProcrServer(server, procrGrpcServer)

	lis, err := net.Listen("tcp", "")
	require.NoError(t, err)

	go server.Serve(lis)
	return server, svc, lis.Addr().String()
}
