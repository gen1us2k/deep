FROM alpine:3.5
MAINTAINER Florin Patan "florinpatan@gmail.com"

WORKDIR /

# DNS stuff
RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

# SSL certs
RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

# Binary
ADD deep /deep

# Runtime
CMD ["/deep"]
