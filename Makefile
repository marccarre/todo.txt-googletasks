.PHONY: lint build test docker-build docker-push clean
.DEFAULT_GOAL := build

IMAGE_USER := marccarre
IMAGE_NAME := todo.txt-googletasks
IMAGE := quay.io/$(IMAGE_USER)/$(IMAGE_NAME)
VERSION := $(shell build/version)

SUPPORTED_GOOS := linux darwin windows

lint:
	docker build --target lint -t $(IMAGE)-lint:latest .

# For each supported operating system, build the binary, and then extract it from the build image:
build:
	@mkdir -p bin
	for os in $(SUPPORTED_GOOS) ; do \
		docker build --target compilation -t $(IMAGE)-build-$$os:latest \
			--build-arg GOOS="$$os" \
			--build-arg VERSION=$(VERSION) \
			. && \
		docker container create --name build-$$os $(IMAGE)-build-$$os:latest && \
		docker container cp build-$$os:/go/src/github.com/marccarre/todo.txt-googletasks/gtasks-$(VERSION)-$$os bin/gtasks-$(VERSION)-$$os && \
		docker container rm -f build-$$os ; \
	done

test:
	rm -f coverage.out
	@docker build --target testing -t $(IMAGE)-testing:latest \
		--build-arg CI=$(CI) \
		--build-arg COVERALLS_TOKEN=$(COVERALLS_TOKEN) \
		--build-arg CLIENT_ID=$(CLIENT_ID) \
		--build-arg CLIENT_SECRET=$(CLIENT_SECRET) \
		--build-arg BASE64_ENCODED_OAUTH_TOKEN=$(BASE64_ENCODED_OAUTH_TOKEN) \
		.
	docker container create --name test $(IMAGE)-testing:latest
	docker container cp test:/go/src/github.com/marccarre/todo.txt-googletasks/coverage.out coverage.out
	docker container rm -f test

docker-build:
	docker build -t $(IMAGE):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		.

docker-push:
	docker push $(IMAGE):$(VERSION)

clean:
	rm -f coverage.out
	rm -fr bin
	-for os in $(SUPPORTED_GOOS) ; do \
		docker container rm -f test ; \
		docker rmi $(IMAGE)-testing:latest ; \
		docker container rm -f build-$$os ; \
		docker rmi $(IMAGE)-build-$$os:latest ; \
		docker rmi $(IMAGE)-lint:latest ; \
	done
	-docker rmi $(IMAGE):$(VERSION)
