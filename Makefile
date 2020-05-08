build:
	go build -ldflags="-w -s -X main.VERSION=$(shell git describe --tags)" .

run:
	go run -ldflags="-w -s -X main.VERSION=$(shell git describe --tags)" .