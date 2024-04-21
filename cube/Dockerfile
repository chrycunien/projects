FROM golang:1.22 as build
WORKDIR /src
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64
COPY go.mod go.sum .
RUN  go mod download
COPY . .
RUN go build -o /bin/cube .

FROM scratch
COPY --from=build /bin/cube /bin/cube
EXPOSE 5555
ENTRYPOINT ["./bin/cube"]