.PHONY: test

build:
	go build cmd/deplist.go

test:
	go test ./... -cover -covermode=atomic
