GOCMD=GO111MODULE=on go

run-server:
	go run *.go server

run-consumer:
	go run *.go consumer

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.coverprofile ./...


bench:
	$(GOCMD) test -bench=. -benchmem ./...

.PHONY: test bench
