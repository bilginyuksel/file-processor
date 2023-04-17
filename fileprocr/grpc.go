package fileprocr

import (
	"bytes"
	"context"

	"github.com/bilginyuksel/file-processor/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedProcrServer

	svc fileprocrService
}

func NewGrpcServer(svc fileprocrService) *GrpcServer {
	return &GrpcServer{svc: svc}
}

func (s *GrpcServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	buf := bytes.NewBuffer(req.GetData())
	if err := s.svc.Store(buf); err != nil {
		return nil, status.Errorf(codes.Internal, "could not process file")
	}
	return &pb.UploadResponse{}, nil
}
