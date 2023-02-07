build:
	go build -o bin/server cmd/main.go


test:
	go test ./... -cover

run: build
	./server

watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make run