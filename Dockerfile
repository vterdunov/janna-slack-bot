FROM golang:1.8.3-alpine AS build-stage

ARG WORKDIR=/go/src/github.com/vterdunov/janna-slack-bot/

WORKDIR $WORKDIR

RUN apk add --no-cache git build-base ca-certificates
RUN go get -v github.com/golang/dep/cmd/dep

COPY . $WORKDIR
RUN [ -d 'vendor' ] || make dep
RUN make compile

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
CMD ["/janna-slack-bot"]
COPY --from=build-stage /go/src/github.com/vterdunov/janna-slack-bot/janna-slack-bot /janna-slack-bot
