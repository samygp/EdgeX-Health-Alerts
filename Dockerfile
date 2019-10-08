FROM golang:1.12.9-alpine3.10 as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN set -x \
  && apk add --update --no-cache \
    jq \
    git \
    gcc \
    libc-dev \
    libgcc \
    make \
    ca-certificates \
  && echo "Installed OS deps."

WORKDIR /go/src/github.com/samygp/edgex-health-alerts
COPY Makefile package.json ./
ENV GO111MODULE on
RUN go mod init
RUN set -x \
  && make install-deps \
  && echo "Installed App deps."

COPY . .
RUN set -x \
  && make static \
  && mv edgex-health-alerts /usr/local/bin/ \
  && echo "Build App."

FROM alpine:3.10

LABEL maintainer="Sam <soysamygp@gmail.com>"

COPY --from=builder /usr/local/bin/edgex-health-alerts /usr/local/bin/
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT ["edgex-health-alerts"]
