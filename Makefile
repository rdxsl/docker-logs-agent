
BINARY := docker-agent-proxy

UNIX_EXECUTABLES := \
	darwin/amd64/$(BINARY) \
	linux/amd64/$(BINARY) \

ALL_EXECUTABLES := $(UNIX_EXECUTABLES:%=bin/%)

GOFILES := $(shell git ls-files | egrep \.go$ | egrep -v ^vendor/ | egrep -v _test.go$)

DOCKER_IMAGE := rdxsl/${BINARY}

DOCKER_REGISTRY := ${DOCKER_REGISTRY}

APP_VERSION ?= $(shell git describe --long)

all: clean test build

build: clean $(ALL_EXECUTABLES)

clean:
	rm -rf bin/ pkg/ $(BINARY)

devcert:
	cd cicd; ./makecert.sh test@test.com

prodcert:
	cd cicd; ./makecert.sh jxie@riotgames.com ..\/conf\/production

# Run unittests
test: docker_image
	bash tests/maketest.sh $(APP_VERSION)

# Run unittests with race condition detector on (takes longer)
testwithrace:
	$(TEST_ENV_VARS) go test $(TEST_FLAGS) -race $(ALL_PACKAGES)


bin/darwin/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo  -ldflags="-w -s" -ldflags="-X main.Version=$(APP_VERSION)" -o "$@" main.go

bin/linux/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo  -ldflags="-w -s" -ldflags="-X main.Version=$(APP_VERSION)" -o "$@" main.go

docker_image: build bin/linux/amd64/$(BINARY) prodcert
	docker build -t $(DOCKER_IMAGE):$(APP_VERSION) -f conf/production/Dockerfile .

docker_debug_image: build bin/linux/amd64/$(BINARY) prodcert
	docker build -t $(DOCKER_IMAGE):$(APP_VERSION) -f conf/production/debug_Dockerfile .

docker_release: docker_image
	docker tag $(DOCKER_IMAGE):$(APP_VERSION) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(APP_VERSION)
