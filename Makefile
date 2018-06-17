PROG_NAME = janna-slack-bot

PORT ?= 8080
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
PROJECT ?= github.com/vterdunov/${PROG_NAME}

GO_VARS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_LDFLAGS :="
GO_LDFLAGS += -s -w
GO_LDFLAGS += -X ${PROJECT}/pkg/version.Commit=${COMMIT}
GO_LDFLAGS += -X ${PROJECT}/pkg/version.BuildTime=${BUILD_TIME}
GO_LDFLAGS +="

TAG ?= $(COMMIT)

all: lint test docker

.PHONY: docker
docker:
	docker build --tag=$(PROG_NAME):$(COMMIT) --tag=$(PROG_NAME):latest --file build/Dockerfile .

.PHONY: push
push:
	docker tag $(PROG_NAME):$(COMMIT) $(PROG_NAME):$(TAG)
	docker push $(PROG_NAME):$(TAG)

.PHONY: dep
dep:
	@dep ensure -v

.PHONY: compile
compile: clean
	 $(GO_VARS) go build -v -ldflags $(GO_LDFLAGS) -o $(PROG_NAME) ./cmd/bot/bot.go

.PHONY: cgo-compile
cgo-compile: clean
	 go build -v -o $(PROG_NAME) ./cmd/bot/bot.go

.PHONY: start
start:
	@env `cat .env | grep -v ^# | xargs` go run ./cmd/bot/bot.go

.PHONY: dc
dc: compile
	docker-compose -f deploy/docker-compose.dev.yml up --build

.PHONY: test
test:
	go test -v ./...

.PHONY: check
check:
	@gometalinter.v2 ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	@rm -f ${PROG_NAME}
