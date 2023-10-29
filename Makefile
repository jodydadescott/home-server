BUILD_NUMBER := latest
PROJECT_NAME := home-dns-server
DOCKER_REGISTRY := jodydadescott
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

default:
	@echo "Make what? linux-amd64 linux-amd64-container darwin-amd64"

linux-amd64:
	@$(MAKE) build/linux-amd64

build/linux-amd64:
	mkdir -p build/linux-amd64
	env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/home-dns-server home-dns-server.go

darwin-amd64:
	@$(MAKE) build/darwin-amd64

build/darwin-amd64:
	mkdir -p build/darwin-amd64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/home-dns-server home-dns-server.go

linux-amd64-container:
	$(MAKE) linux-amd64
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

push:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

clean:
	$(RM) -rf build
