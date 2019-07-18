FROM fluent/fluentd-kubernetes-daemonset:v1.4-debian-elasticsearch-1

LABEL author="smanpathak@platform9.com"

USER root
WORKDIR /home/fluent
ENV PATH /fluentd/vendor/bundle/ruby/2.6.0/bin:$PATH
ENV GEM_PATH /fluentd/vendor/bundle/ruby/2.6.0
ENV GEM_HOME /fluentd/vendor/bundle/ruby/2.6.0
# skip runtime bundler installation
ENV FLUENTD_DISABLE_BUNDLER_INJECTION 1

RUN gem install fluent-plugin-s3
