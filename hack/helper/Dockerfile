FROM alpine:latest
LABEL author="smanpathak@platform9.com"

RUN mkdir -p /helper/bin
WORKDIR /helper
COPY build/bin/fluentd-operator-helper bin/helper
RUN chmod +x bin/helper

ENTRYPOINT [ "bin/helper" ]