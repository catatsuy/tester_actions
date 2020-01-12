export GO111MODULE=on

all: bin/bento

bin/bento: cmd/bento/main.go
	go build -o bin/bento cmd/bento/main.go
