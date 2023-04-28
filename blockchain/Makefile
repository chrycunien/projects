BINARY_NAME=blockchain
DB_NAME=blockchain.db

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean_binary:
	rm ${BINARY_NAME}

clean_db:
	rm ${DB_NAME}

clean: clean_binary clean_db	

clean_all:
	go clean
	rm ${BINARY_NAME}*

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all