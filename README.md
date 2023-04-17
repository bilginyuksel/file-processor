# file-processor

Application allows clients to upload file either using REST API or gRPC API. File will be written to the disk by chunks so that allows clients to upload unlimited size files without crushing the app memory.

## Getting Started

Run the application

```bash
make up
```

Run tests

```bash
make test
```

Create auto-generated codes

```bash
make gen
```

## Example

After running the testcases to test the application manually, you can use the guidance below.

Test the application using REST endpoint

```bash
curl --location 'localhost:8010/files' --form 'file=@"<absolute-file-path>"'
```

Test the application using gRPC endpoint or use gRPC curl to upload the file

```bash
go run examples/grpc_client.go
```
