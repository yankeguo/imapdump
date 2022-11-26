FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /imapdump ./cmd/imapdump

FROM ghcr.io/guoyk93/acicn/alpine:3.16

WORKDIR /data

COPY --from=builder /imapdump /imapdump

ENV MINIT_MAIN          /imapdump
ENV MINIT_MAIN_DIR      /data
ENV MINIT_MAIN_NAME     imapdump
ENV MINIT_MAIN_KIND     cron
ENV MINIT_MAIN_CRON     "@every 6h"