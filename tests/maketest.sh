#!/usr/bin/env bash
if [ "$1" == "" ]; then
  echo "need a docker version to test"
  exit 1
fi
VERSION=$1

docker run -d -v /var/run/docker.sock:/var/run/docker.sock  -p 7001:7001 -p 7002:7002 \
   --name docker-agent-proxy-test-${VERSION} rdxsl/docker-agent-proxy:${VERSION}

statusCode=$(curl --write-out %{http_code} --silent --output /dev/null -X GET "http://localhost:7001/v1/containers/docker-agent-proxy-test-${VERSION}/logs/?tail=5" -H  "accept: application/json")
if [ $statusCode != 200 ]; then
  echo "Can't run the log test to http port 7001, please check! status=$statusCode"
  exit 1
else
  echo "log test to http port 7001 PASSED."
fi

statusCode=$(curl --write-out %{http_code} --silent --output /dev/null --key conf/production/certs/client.key --cert conf/production/certs/client.pem -k -X GET "https://localhost:7002/v1/containers/docker-agent-proxy-test-${VERSION}/logs/?tail=5" -H  "accept: application/json")
if [ $statusCode != 200 ]; then
  echo "Can't run the log test to https port 7002, please check! status=$statusCode"
  exit 1
else
  echo "log test to https port 7002 PASSED."
fi

statusCode=$(curl --write-out %{http_code} --silent --output /dev/null  --key conf/production/certs/client.key --cert conf/production/certs/client.pem -k -i -X POST "https://localhost:7002/v1/containers/docker-agent-proxy-test-${VERSION}/exec" -H  "accept: application/json" -H  "content-type: application/json" -d "{\"Cmd\":[\"ls\",\"/bin\"]}"   )
if [ $statusCode != 200 ]; then
  echo "Can't run the exec test to https port 7002, please check! status=$statusCode"
  exit 1
else
  echo "exec test to https port 7002 PASSED."
fi

####### FAIL TEST FOR BAD CERT #######
echo "bad key" > conf/production/certs/bad-client.key
statusCode=$(curl --write-out %{http_code} --silent --output /dev/null --key conf/production/certs/bad-client.key --cert conf/production/certs/client.pem -k -X GET "https://localhost:7002/v1/containers/docker-agent-proxy-test-${VERSION}/logs/?tail=5" -H  "accept: application/json")
if [ $statusCode != 000 ]; then
  echo "Can't run the best cert log test to https port 7002, please check! status=$statusCode"
  exit 1
else
  echo "BAD cert log test to https port 7002 PASSED."
fi

docker rm -f docker-agent-proxy-test-${VERSION}
