.PHONY=run main

run:
	go run *.go

build:
	go build -o iap *.go
