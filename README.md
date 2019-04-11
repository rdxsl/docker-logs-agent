# docker-logs-agent

Web wrapper to access docker logs via unix socket `/var/run/docker.sock`. Only use this in a secure network environment.

## Testing
```
make docker_image

docker run --privileged -v /var/run/docker.sock:/var/run/docker.sock  -p 7001:7001 rdxsl/docker-logs-agent:$VERSION

curl -X GET "http://localhost:7001/v1/containers/$containerID/logs"
```

## Deploy
```
export DOCKER_REGISTRY=your_docker_reigsrey

docker login your_docker_reigsrey

make docker_release
``
