.PHONY: test

build:
	go build cmd/deplist/deplist.go

test:
	go test ./... -cover -covermode=atomic
