.PHONY: run
run:
	@go build -buildmode=plugin -o ./app1/app1.so ./app1/app1.go
	@go build -buildmode=plugin -o ./app2/app2.so ./app2/app2.go
	@go run ./ ../../pride-and-prejudice.txt

.PHONY: test
test:
	@go test -v -count=1 ./...

.PHONY: bench
bench:
	@go test -bench=.