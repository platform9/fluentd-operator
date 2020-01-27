FROM alpine:latest
LABEL author="smanpathak@platform9.com"

RUN mkdir -p /fluentd/bin
RUN apk add --no-cache \
        libc6-compat
WORKDIR /fluentd
COPY build/bin/fluentd-operator-linux-amd64 bin/fluentd-operator
ADD etc etc
RUN chmod +x bin/fluentd-operator

ENTRYPOINT [ "bin/fluentd-operator" ]
