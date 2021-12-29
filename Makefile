.PHONY: init build tidy test bash
GO_VERSION := 1.16.4
PLUGIN_NAME := cached-router

init:
	docker run -v "${PWD}:/plugin" golang:$(GO_VERSION) bash -c "cd /plugin && go mod init github.com/RePrete/$(PLUGIN_NAME)" 

build:
	docker run -v "${PWD}:/plugin" golang:$(GO_VERSION) bash -c "cd /plugin && go build --buildmode=plugin -o ./build/$(PLUGIN_NAME).so ./" 

tidy:
	docker run -v "${PWD}:/plugin" golang:$(GO_VERSION) bash -c "cd /plugin && go mod tidy" 

test:
	docker run -v "${PWD}:/plugin" golang:$(GO_VERSION) bash -c "cd /plugin && go test -v ./..." 

bash:
	docker run -v "${PWD}:/plugin" -it golang:$(GO_VERSION) bash