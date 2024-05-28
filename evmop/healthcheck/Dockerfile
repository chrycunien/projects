FROM golang:1.22 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o server main.go
EXPOSE 8080
ENV PORT=8080

FROM scratch
WORKDIR /
COPY --from=builder /workspace/server .
USER 65532:65532

ENTRYPOINT [ "/server" ]