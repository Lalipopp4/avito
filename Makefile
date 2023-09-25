all:
	go build -o app ./cmd/app
	./app

build:
	go build -o app ./cmd/app

run:
	./app