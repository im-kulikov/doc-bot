# Builder
FROM golang:1.12-alpine3.9 as builder

COPY . /src/doc-bot

WORKDIR /src/doc-bot

ARG VERSION=0.0.0-unknown
ARG REPO=none

RUN set -x \
    && export BUILD=$(date -u +%s%N) \
    && export LDFLAGS="-w -s -X ${REPO}/misc.Version=${VERSION} -X ${REPO}/misc.Build=${BUILD}" \
    && export CGO_ENABLED=0 \
    && go build -mod=vendor -v -ldflags "${LDFLAGS}" -o /go/bin/docbot ./cmd/bot \
    && chmod 1755 /go/bin/docbot

# Executable image
FROM scratch

ENV BOT_TELEGRAM_TOKEN=""
ENV BOT_TELEGRAM_PROXY=""

WORKDIR /

COPY --from=builder /go/bin/docbot /docbot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/doc-bot/config.yml /config.yml

CMD ["/docbot"]
