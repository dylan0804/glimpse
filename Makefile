APP_NAME=glimpse

.PHONY: build run test

build: 
	wails build

run: 
	wails dev

test:
	go test -v ./screenshots