<<<<<<< HEAD
TEST?=./...

default: test

fmt: generate
	go fmt ./...

test: generate
	go get -t ./...
	go test $(TEST) $(TESTARGS)

generate:
	go generate ./...

updatedeps:
	go get -u golang.org/x/tools/cmd/stringer

.PHONY: default generate test updatedeps
=======
TEST?=./...

default: test

fmt: generate
	go fmt ./...

test: generate
	go get -t ./...
	go test $(TEST) $(TESTARGS)

generate:
	go generate ./...

updatedeps:
	go get -u golang.org/x/tools/cmd/stringer

.PHONY: default generate test updatedeps
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2
