export GO111MODULE=on

all: bin/bento

bin/bento: cmd/bento/main.go mirait/*.go
	go build -o bin/bento cmd/bento/main.go

vet:
	go vet ./...

errcheck:
	errcheck ./...

staticcheck:
	staticcheck -checks="all,-ST1000" ./...

clean:
	rm -rf bin/*
