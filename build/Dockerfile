ARG GOLANG_VERSION=1.11.2

# ---------------
FROM golang:${GOLANG_VERSION} AS env

RUN mkdir /user && \
    echo 'janna-bot:x:65534:65534:janna-bot:/:' > /user/passwd && \
    echo 'janna-bot:x:65534:' > /user/group

WORKDIR /gomod
COPY go.mod go.sum ./
RUN go mod download

# ---------------
FROM env AS builder
WORKDIR /build
COPY . .

RUN make compile

# ---------------
FROM builder AS tester

RUN make test

# ---------------
FROM alpine:3.8 AS final

ENV PORT=8080

WORKDIR /janna-bot

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER janna-bot:janna-bot

CMD ["./janna-bot"]

COPY --from=builder /build/janna-bot /janna-bot
