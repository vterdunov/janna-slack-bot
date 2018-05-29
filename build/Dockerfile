FROM golang:1.10.2-alpine AS build-stage

WORKDIR /go/src/github.com/vterdunov/janna-slack-bot/

RUN apk add --no-cache git build-base ca-certificates

COPY . .

RUN make compile

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
CMD ["/janna-slack-bot"]

COPY --from=build-stage /go/src/github.com/vterdunov/janna-slack-bot/janna-slack-bot /janna-slack-bot
