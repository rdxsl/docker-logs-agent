# start from alpine
FROM alpine
MAINTAINER Jack Xie <jxie@riotgames.com>

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY bin/linux/amd64/docker-logs-agent /go/bin/docker-logs-agent
RUN chown guest:users /go/bin/docker-logs-agent

COPY conf/production/app_prod.conf /go/bin/conf/app.conf

EXPOSE 7000

ENTRYPOINT ["/go/bin/docker-logs-agent"]
