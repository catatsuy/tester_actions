export GO111MODULE=on

all: bin/bento

bin/bento: cmd/bento/main.go mirait/*.go config/*.go cli/*.go
	go mod tidy
	go build -ldflags "-X github.com/catatsuy/bento/cli.Version=`git rev-list HEAD -n1`" -o bin/bento cmd/bento/main.go

vet:
	go vet ./...

errcheck:
	errcheck ./...

staticcheck:
	staticcheck -checks="all,-ST1000" ./...

clean:
	rm -rf bin/*

check:
	go test ./...

.PHONY: all vet errcheck staticcheck clean check
