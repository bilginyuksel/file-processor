package fileprocr

import (
	"context"

	"github.com/bilginyuksel/file-processor/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedProcrServer
}

func (s *GrpcServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
