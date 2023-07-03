FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /imapdump ./cmd/imapdump

FROM guoyk/minit:1.11.1 AS minit

FROM alpine:3.18

# install packages
RUN apk add --no-cache tzdata ca-certificates

# install minit
RUN mkdir -p /opt/bin
ENV PATH "/opt/bin:${PATH}"
COPY --from=minit /minit /opt/bin/minit
ENV MINIT_LOG_DIR none
ENTRYPOINT ["/opt/bin/minit"]

WORKDIR /data

COPY --from=builder /imapdump /imapdump

ENV MINIT_MAIN          /imapdump
ENV MINIT_MAIN_DIR      /data
ENV MINIT_MAIN_NAME     imapdump
ENV MINIT_MAIN_KIND     cron
ENV MINIT_MAIN_CRON     "@every 6h"