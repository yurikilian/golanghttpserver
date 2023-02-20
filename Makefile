build:
	go build -o bin/server cmd/main.go


test:
	go test ./... -cover

run: build
	GOMEMLIMIT=256MiB GOMAXPROCS=8 ./bin/server

watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make run