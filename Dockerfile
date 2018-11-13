# ------------------------------------------------------------------------ setup
FROM golang:1.11.2-alpine3.8 AS setup

# Install git, as not present in golang:1.11.2-alpine3.8 and required by dep ensure -vendor-only
RUN apk --no-cache add git
# Install bash & curl, as not present in golang:1.11.2-alpine3.8 and required by codecov
RUN apk --no-cache add bash curl

# Install dep, for dependencies management in Go:
RUN curl -fsSLo dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
	echo "287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep" | sha256sum -c && \
	chmod +x dep && \
	mv dep $GOPATH/bin/dep

# Install gometalinter and all underlying linters, for code quality:
RUN curl -fsSLo gometalinter.tar.gz https://github.com/alecthomas/gometalinter/releases/download/v2.0.11/gometalinter-2.0.11-linux-amd64.tar.gz && \
	echo "97d8bd0a4d024740964c7fc2ae41276cf5f839ccf0749528ca900942f656d201  gometalinter.tar.gz" | sha256sum -c && \
	tar xzf gometalinter.tar.gz && \
	cd gometalinter-2.0.11-linux-amd64 && \
	rm -f COPYING README.md && \
	chmod +x * && \
	mv * $GOPATH/bin && \
	cd .. && \
	rm -rf gometalinter-2.0.11-linux-amd64

# Install markdownlint and its dependencies:
RUN apk --no-cache add ruby ruby-json && \
	gem install --no-ri --no-rdoc mdl && \
	ruby -v && gem -v && mdl --version

# Install goveralls, for code coverage:
RUN go get github.com/mattn/goveralls

# Gopkg.toml and Gopkg.lock lists project dependencies.
# These layers will only be re-built when Gopkg files are updated:
COPY Gopkg.lock Gopkg.toml /go/src/github.com/marccarre/todo.txt-googletasks/
WORKDIR /go/src/github.com/marccarre/todo.txt-googletasks
# Install all dependencies:
RUN dep ensure -vendor-only

# ------------------------------------------------------------------------- lint
FROM setup AS lint

# Copy project. This layer will be rebuilt when ever a file has changed in the project directory
COPY . /go/src/github.com/marccarre/todo.txt-googletasks

RUN find . -name "*.md" -not -path "./vendor/*" -exec mdl --rules ~MD013 {} \; && \
	gometalinter $(go list ./...)

# ------------------------------------------------------------------------ build
FROM setup AS build

# Copy project. This layer will be rebuilt when ever a file has changed in the project directory
COPY . /go/src/github.com/marccarre/todo.txt-googletasks

# Set the provided GOOS, or default it to "linux":
ARG GOOS=linux
# Set the provided VERSION:
ARG VERSION

# Compile for the configured operating system:
# -tags netgo -ldflags: use the built-in net package
# -w: disable debug information for smaller binary
# -extldflags "-static": build a static binary to avoid having to install 3rd party libraries
# -X pkg/version.Version=$(VERSION): injects this binary's version, so it can be printed later on
RUN CGO_ENABLED=0 GOARCH=amd64 go build \
	-tags netgo -ldflags \
	"-w -extldflags '-static' -X github.com/marccarre/todo.txt-googletasks/pkg/version.Version=${VERSION}" \
	-o gtasks-${VERSION}-${GOOS} cmd/gtasks/gtasks.go

# ------------------------------------------------------------------------- test
FROM build AS test

# Set the provided Google Tasks API credentials, as well as the OAuth token,
# base64-encoded, as we cannot authenticate online within the container.
ARG CLIENT_ID
ARG CLIENT_SECRET
ENV OAUTH_SCOPES=https://www.googleapis.com/auth/tasks
ARG BASE64_ENCODED_OAUTH_TOKEN
RUN echo -n "${BASE64_ENCODED_OAUTH_TOKEN}" | base64 -d > \
	~/.cache/todo.txt-googletasks_"$(echo -n "${CLIENT_ID}${CLIENT_SECRET}${OAUTH_SCOPES}" | sha256sum | awk '{ print $1 }' | head -c 8)"

# Set the provided COVERALLS_TOKEN, or default it to empty string otherwise:
ARG COVERALLS_TOKEN
# Set the provided CODECOV_TOKEN, or default it to empty string otherwise:
ARG CODECOV_TOKEN
# Set the provided CI env. var., or default it to empty string otherwise:
ARG CI

# Run tests and, optionally, upload code coverage to coveralls.io and codecov.io:
RUN CGO_ENABLED=0 go test -v -timeout 30s -cover -covermode=count -coverprofile=coverage.out ./...
RUN if [ "$CI" == "true" ] && [ ! -z "$COVERALLS_TOKEN" ] && [ ! -z "$CODECOV_TOKEN" ]; then \
	goveralls \
	-coverprofile=/go/src/github.com/marccarre/todo.txt-googletasks/coverage.out \
	-service=circle-ci -repotoken=$COVERALLS_TOKEN ; \
	curl -fsSLo codecov.sh https://codecov.io/bash ; chmod +x codecov.sh ; ./codecov.sh -t $CODECOV_TOKEN ; else \
	echo "Skipped upload of code coverage." ; fi


# ---------------------------------------------------------------------- runtime
FROM scratch
ARG VERSION
COPY --from=build /go/src/github.com/marccarre/todo.txt-googletasks/gtasks-${VERSION}-linux /gtasks
ENTRYPOINT ["/gtasks"]
CMD ["--help"]
