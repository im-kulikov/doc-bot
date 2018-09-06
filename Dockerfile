# Builder
FROM golang:1.11.0-alpine3.8 as builder

COPY . /go/src/github.com/im-kulikov/doc-bot

WORKDIR /go/src/github.com/im-kulikov/doc-bot

RUN set -x \
    && apk add --no-cache git \
    && export VERSION=$(git rev-parse --verify HEAD) \
    && export BUILD=$(date -u +%s%N) \
    && export LDFLAGS="-w -s -X main.BuildVersion=${VERSION} -X main.BuildTime=${BUILD}" \
    && export CGO_ENABLED=0 \
    && go build -v -ldflags "${LDFLAGS}" -o /go/bin/docbot . \
    && chmod 1755 /go/bin/docbot

# Executable image
FROM scratch

ENV BOT_TELEGRAM_TOKEN=""
ENV BOT_TELEGRAM_PROXY=""

WORKDIR /

COPY --from=builder /go/bin/docbot /docbot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/docbot"]