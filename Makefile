PROG_NAME = janna-bot
IMAGE_NAME = vterdunov/$(PROG_NAME)

PORT ?= 8080
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
PROJECT ?= github.com/vterdunov/${PROG_NAME}

GO_VARS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_LDFLAGS :="
GO_LDFLAGS += -s -w
GO_LDFLAGS += -X ${PROJECT}/internal/version.Commit=${COMMIT}
GO_LDFLAGS += -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}
GO_LDFLAGS +="

TAG ?= $(COMMIT)

GOLANGCI_LINTER_VERSION = v1.12.2

all: lint docker

.PHONY: docker
docker: ### Build Docker container
	docker build --tag=$(IMAGE_NAME):$(COMMIT) --tag=$(IMAGE_NAME):latest --file build/Dockerfile .

.PHONY: push
push: ### Push docker container to registry
	docker tag $(IMAGE_NAME):$(COMMIT) $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME):$(TAG)

.PHONY: compile
compile: clean ### Compile bot
	$(GO_VARS) go build -v -ldflags $(GO_LDFLAGS) -o $(PROG_NAME) ./cmd/bot/bot.go

.PHONY: cgo-compile
cgo-compile: clean
	go build -v -o $(PROG_NAME) ./cmd/bot/bot.go

.PHONY: run
run: ### Extract env variables from .env and run bot with race detector
	@env `cat .env | grep -v ^# | xargs` go run -race ./cmd/bot/bot.go

.PHONY: run-docker
run-docker: docker ### Build docker image. Extract env variables from .env and run bot using docker.
	docker run -it --rm --env-file=.env $(IMAGE_NAME):latest

.PHONY: test
test: ### Run tests
	go test -v ./...

.PHONY: lint
lint: ### Run linters
	@echo Linting...
	@docker run -it --rm -v $(CURDIR):/lint -w /lint golangci/golangci-lint:$(GOLANGCI_LINTER_VERSION) golangci-lint run

.PHONY: clean
clean:
	@rm -f ${PROG_NAME}

.PHONY: help
help: ### Show this help.
	@sed -e '/__hidethis__/d; /###/!d; s/:.\+### /\t/g' $(MAKEFILE_LIST)
