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

- cli
```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

- helm
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-nginx bitnami/nginx
POD_NAME=$(kubectl get pod \
    --selector=app.kubernetes.io/name=nginx \
    --template '{{index .items 0 "metadata" "name" }}')
SERVICE_IP=$(kubectl get svc \
    --namespace default my-nginx --template "{{ .spec.clusterIP }}")
helm uninstall my-nginx
```
```bash
mkdir deploy && cd deploy
helm create proglog
helm template proglog
rm proglog/templates/**/*.yaml proglog/templates/NOTES.txt

helm install proglog deploy/proglog
helm uninstall proglog
```
