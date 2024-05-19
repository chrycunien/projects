# README

## Requirements

- protobuf
```bash
# protoc
brew install protobuf

# protobuf and grpc runtime
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# dependencies for generated code
go get google.golang.org/grpc@latest
```