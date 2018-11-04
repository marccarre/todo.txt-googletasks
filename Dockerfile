# ------------------------------------------------------------------------ setup
FROM golang:1.11.2-alpine3.8 AS setup

# Install git, as not present in golang:1.11.2-alpine3.8 and required by dep ensure -vendor-only
RUN apk --no-cache add git

# Install dep, for dependencies management in Go:
RUN apk --no-cache add curl && \
	curl -fsSLo dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
	echo "287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep" | sha256sum -c && \
	chmod +x dep && \
	mv dep $GOPATH/bin/dep && \
	apk del curl && \
	rm -rf /var/cache/apk/*

# Install goveralls, for code coverage:
RUN go get github.com/mattn/goveralls

# Gopkg.toml and Gopkg.lock lists project dependencies.
# These layers will only be re-built when Gopkg files are updated:
COPY Gopkg.lock Gopkg.toml /go/src/github.com/marccarre/todo.txt-googletasks/
WORKDIR /go/src/github.com/marccarre/todo.txt-googletasks
# Install all dependencies:
RUN dep ensure -vendor-only

# ------------------------------------------------------------------ compilation
FROM setup AS compilation

# Copy project. This layer will be rebuilt when ever a file has changed in the project directory
COPY . /go/src/github.com/marccarre/todo.txt-googletasks

# Set the provided GOOS, or default it to "linux":
ARG GOOS=linux
ENV GOOS=$GOOS

# Compile for the configured operating system:
# -tags netgo -ldflags: use the built-in net package
# -w: disable debug information for smaller binary
# -extldflags "-static": build a static binary to avoid having to install 3rd party libraries
RUN CGO_ENABLED=0 GOARCH=amd64 go build \
	-tags netgo -ldflags \
	'-w -extldflags "-static"' \
	-o gtasks-${GOOS} cmd/gtasks/gtasks.go

# ---------------------------------------------------------------------- testing
FROM compilation AS testing

# Set the provided Google Tasks API credentials:
ARG CLIENT_ID
ENV CLIENT_ID=$CLIENT_ID
ARG CLIENT_SECRET
ENV CLIENT_SECRET=$CLIENT_SECRET
# Set the provided COVERALLS_TOKEN, or default it to empty string otherwise:
ARG COVERALLS_TOKEN
ENV COVERALLS_TOKEN=$COVERALLS_TOKEN
# Set the provided CI env. var., or default it to empty string otherwise:
ARG CI
ENV CI=$CI

# Run tests and, optionally, upload code coverage to coveralls.io:
RUN CGO_ENABLED=0 go test -v -timeout 30s -cover -covermode=count -coverprofile=coverage.out ./...
RUN [ "$CI" == "true" ] && [ ! -z "$COVERALLS_TOKEN" ] && \
	goveralls \
	-coverprofile=/go/src/github.com/marccarre/todo.txt-googletasks/coverage.out \
	-service=circle-ci \
	-repotoken=$COVERALLS_TOKEN \
	|| true

# ---------------------------------------------------------------------- runtime
FROM scratch
COPY --from=compilation /go/src/github.com/marccarre/todo.txt-googletasks/gtasks-linux /gtasks
ENTRYPOINT ["/gtasks"]
CMD ["--help"]
