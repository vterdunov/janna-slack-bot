FROM golang:1.8.3-alpine AS build-stage

WORKDIR /build
RUN apk update && \
    apk add --no-cache git build-base
RUN go get -v github.com/adampointer/go-slackbot
COPY main.go /build
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-s -w" -o janna-slack-bot main.go

FROM scratch
CMD ["/janna-slack-bot"]
COPY --from=build-stage /build/janna-slack-bot /janna-slack-bot
