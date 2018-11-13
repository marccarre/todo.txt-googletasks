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
		docker build --target build -t $(IMAGE)-build-$$os:latest \
			--build-arg GOOS="$$os" \
			--build-arg VERSION=$(VERSION) \
			. && \
		docker container create --name build-$$os $(IMAGE)-build-$$os:latest && \
		docker container cp build-$$os:/go/src/github.com/marccarre/todo.txt-googletasks/gtasks-$(VERSION)-$$os bin/gtasks-$(VERSION)-$$os && \
		docker container rm -f build-$$os ; \
	done

test:
	@docker build --target test -t $(IMAGE)-test:latest \
		--build-arg CI=$(CI) \
		--build-arg COVERALLS_TOKEN=$(COVERALLS_TOKEN) \
		--build-arg CODECOV_TOKEN=$(CODECOV_TOKEN) \
		--build-arg CLIENT_ID=$(CLIENT_ID) \
		--build-arg CLIENT_SECRET=$(CLIENT_SECRET) \
		--build-arg BASE64_ENCODED_OAUTH_TOKEN=$(BASE64_ENCODED_OAUTH_TOKEN) \
		.

docker-build:
	docker build -t $(IMAGE):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		.

docker-push:
	docker push $(IMAGE):$(VERSION)

clean:
	rm -fr bin
	-docker rmi $(IMAGE):$(VERSION)
	-docker rmi $(IMAGE)-test:latest
	-for os in $(SUPPORTED_GOOS) ; do \
		docker container rm -f build-$$os ; \
		docker rmi $(IMAGE)-build-$$os:latest ; \
	done
	-docker rmi $(IMAGE)-lint:latest
