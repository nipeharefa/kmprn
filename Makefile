GOCMD=GO111MODULE=on go

run-server:
	go run *.go server

run-consumer:
	go run *.go consumer

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.coverprofile ./...

lint:
	golint -set_exit_status ./...

bench:
	$(GOCMD) test -bench=. -benchmem ./...

.PHONY: test bench
