PROG_NAME = janna-slack-bot
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GO_VARS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_LDFLAGS=-v -ldflags="-s -w"

start: build run

build:
	docker build -t $(PROG_NAME) .

run:
	docker run -it --rm --name=$(PROG_NAME) -p 4567:4567 --env-file=.env $(PROG_NAME)

dep:
	@dep ensure -v

compile:
	 $(GO_VARS) go build $(GO_LDFLAGS) -o $(PROG_NAME)
