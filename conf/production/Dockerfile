# start from alpine
FROM alpine
MAINTAINER Jack Xie <jxie@riotgames.com>

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY bin/linux/amd64/docker-agent-proxy /go/bin/docker-agent-proxy
RUN chown guest:users /go/bin/docker-agent-proxy

COPY conf/production/app_prod.conf /go/bin/conf/app.conf

COPY conf/production/certs /go/bin/conf/certs

EXPOSE 7001:7002

ENTRYPOINT ["/go/bin/docker-agent-proxy"]
