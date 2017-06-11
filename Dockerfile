FROM golang:1.8.3-alpine AS build-stage

WORKDIR /build
RUN apk add --no-cache git build-base ca-certificates
RUN go get -v github.com/adampointer/go-slackbot
COPY main.go /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o janna-slack-bot main.go

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
CMD ["/janna-slack-bot"]
COPY --from=build-stage /build/janna-slack-bot /janna-slack-bot
