package pb

// Generates grpc server and client code.
//go:generate protoc --proto_path=. --go_out=. --go_opt=module=github.com/bilginyuksel/file-processor/pb fileprocr.proto
//go:generate protoc --proto_path=. --go-grpc_out=. --go-grpc_opt=module=github.com/bilginyuksel/file-processor/pb fileprocr.proto
