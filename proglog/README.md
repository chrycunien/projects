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
go get github.com/grpc-ecosystem/go-grpc-middleware
```

- TLS certificates
```bash
go install github.com/cloudflare/cfssl/cmd/cfssl@latest

go install github.com/cloudflare/cfssl/cmd/cfssljson@latest
```

- ACL
```bash
go get github.com/casbin/casbin
```

- Observability
```bash
go get go.uber.org/zap
go get go.opencensus.io
go get github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
```

- Discovery
```bash
go get github.com/hashicorp/serf/serf
go get github.com/travisjeffery/go-dynaport
```

- Consensus
```bash
go get github.com/hashicorp/raft
go mod edit -replace github.com/hashicorp/raft-boltdb=github.com/travisjeffery/raft-boltdb@v1.0.0
go get github.com/hashicorp/raft-boltdb
go get github.com/soheilhy/cmux
```