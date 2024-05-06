# README

## Install gRPC
```bash
# protoc
brew install protobuf

# protoc -- go protobuf message plugin
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# protoc -- go protobuf service plugin
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# go packages for generated code
go get google.golang.org/grpc@latest
```

## Generate Stub
```bash
protoc -I ./order \
    --go_out ./golang \
    --go_opt paths=source_relative \
    --go-grpc_out ./golang \
    --go-grpc_opt paths=source_relative \
    ./order/order.proto

# to preserve order structure => ./golang/order/*
protoc -I . \
    --go_out ./golang \
    --go_opt paths=source_relative \
    --go-grpc_out ./golang \
    --go-grpc_opt paths=source_relative \
    ./order/order.proto

# flattened structure => ./golang/*
protoc -I ./order \
    --go_out ./golang \
    --go_opt paths=source_relative \
    --go-grpc_out ./golang \
    --go-grpc_opt paths=source_relative \
    ./order/order.proto

# cheatsheet
protoc -I <import_dir (default=.)> \
    --go_out <out_dir> \
    --go_opt paths=source_relative \
    --go-grpc_out <out_dir> \
    --go-grpc_opt paths=source_relative \
    /path/to/generate/*.proto /path2/to/generate/*.proto ...

protoc -I . \
    --go_out ./golang \
    --go_opt paths=source_relative \
    --go-grpc_out ./golang \
    --go-grpc_opt paths=source_relative \
    ./order/order.proto ./payment/payment.proto ./shipping/shipping.proto
```

## Tagging
```bash
git tag -a golang/order/v1.2.3 -m "golang/order/v1.2.3"
git tag -a golang/payment/v1.2.8 -m "golang/payment/v1.2.8"
git tag -a golang/shipping/v1.2.6 -m "golang/shipping/v1.2.6"
git push --tags

go get -u microgrpc-proto/golang/order@latest
go get -u microgrpc-proto/golang/order@v1.2.3
```

## Reference
- https://stackoverflow.com/questions/60578892/protoc-gen-go-grpc-program-not-found-or-is-not-executable
- https://stackoverflow.com/questions/70731053/protoc-go-opt-paths-source-relative-vs-go-grpc-opt-paths-source-relative
- 