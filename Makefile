.PHONY: build docker-build docker-push clean
.DEFAULT_GOAL := build

IMAGE_USER := marccarre
IMAGE_NAME := todo.txt-googletasks
IMAGE := $(IMAGE_USER)/$(IMAGE_NAME)

SUPPORTED_GOOS := linux darwin windows

# For each supported operating system, build the binary, and then extract it from the build image:
build:
	@mkdir -p bin
	for os in $(SUPPORTED_GOOS) ; do \
		docker build --target compilation -t $(IMAGE)-build-$$os:latest --build-arg GOOS="$$os" . && \
		docker container create --name build-$$os $(IMAGE)-build-$$os:latest && \
		docker container cp build-$$os:/go/src/github.com/marccarre/todo.txt-googletasks/gtasks-$$os bin/gtasks-$$os && \
		docker container rm -f build-$$os ; \
	done

docker-build:
	docker build -t quay.io/$(IMAGE):latest .

docker-push:
	docker push quay.io/$(IMAGE):latest

clean:
	rm -fr bin
	-for os in $(SUPPORTED_GOOS) ; do \
		docker container rm -f build-$$os ; \
		docker rmi $(IMAGE)-build-$$os:latest ; \
	done
	-docker rmi $(IMAGE):latest
