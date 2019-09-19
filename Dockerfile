FROM centos:centos7.4.1708

COPY shortlink /usr/bin/shortlink
RUN mkdir -p /config
COPY config /config

ENTRYPOINT ["/usr/bin/shortlink"]
