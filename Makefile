build:
	go build -o bin/server cmd/main.go

run: build
	./server

watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make run