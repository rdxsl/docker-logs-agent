# docker-agent-proxy

Web wrapper to access docker logs via unix socket `/var/run/docker.sock`. Secured via client side tls certs.

## Versioning
run the following command to add new version
```
git tag 1-X -m "add some message"
```

## Testing
```
make test
```

## Production docker images
```
generate certs using cicd/makecert.sh

if you already have saved certs: copy certs/* conf/production/certs/

make docker_image
```

## Deploy
```
export DOCKER_REGISTRY=your_docker_reigsrey

docker login your_docker_reigsrey

make docker_release
``
