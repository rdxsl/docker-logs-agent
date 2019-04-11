
BINARY := docker-logs-agent

UNIX_EXECUTABLES := \
	darwin/amd64/$(BINARY) \
	linux/amd64/$(BINARY) \

ALL_EXECUTABLES := $(UNIX_EXECUTABLES:%=bin/%)

GOFILES := $(shell git ls-files | egrep \.go$ | egrep -v ^vendor/ | egrep -v _test.go$)

DOCKER_IMAGE := rdxsl/docker-logs-agent

DOCKER_REGISTRY := ${DOCKER_REGISTRY}

APP_VERSION ?= $(shell cat version/version.go | grep 'Version = ' | cut -d ' ' -f 4 | sed 's/"//g')

all: clean test build

build: clean $(ALL_EXECUTABLES)

clean:
	rm -rf bin/ pkg/ $(BINARY)

# Run unittests
test:
	$(TEST_ENV_VARS) go test $(TEST_FLAGS) $(ALL_PACKAGES)

# Run unittests with race condition detector on (takes longer)
testwithrace:
	$(TEST_ENV_VARS) go test $(TEST_FLAGS) -race $(ALL_PACKAGES)


bin/darwin/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags="-s $(VERSION_UPDATE_FLAG)" -o "$@" main.go

bin/linux/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-s $(VERSION_UPDATE_FLAG)" -o "$@" main.go

docker_image: build bin/linux/amd64/$(BINARY)
	ls -al bin/linux/amd64/$(EXECUTABLE)
	docker build -t $(DOCKER_IMAGE):$(APP_VERSION) -f conf/production/Dockerfile .

docker_release: docker_image
	docker tag $(DOCKER_IMAGE):$(APP_VERSION) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
