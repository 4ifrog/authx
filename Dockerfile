#####################
# --- Build Stage ---
#####################

# Pin version as much as possible
FROM golang:1.14-alpine3.12 as build

# Install system dependencies for the build
RUN apk add --no-cache \
    ca-certificates \
    git \
    make

# Set the build environments
WORKDIR /go/src

# Install build dependencies
COPY Makefile ./
COPY go.mod go.sum ./
RUN make install

# Copy the project and build
COPY . ./
RUN mkdir -p bin
RUN make build

#####################
# --- Final Stage ---
#####################

FROM alpine:3.12

# Install tini - need for production for graceful shutdowns
RUN apk add --no-cache tini

# Install system runtime dependencies
RUN apk add --no-cache \
    ca-certificates

# Copy binaries and config file over
WORKDIR /go/bin
USER nobody:nobody
COPY --from=build /go/src/bin/authx ./
COPY --from=build /go/src/bin/config.yaml ./

# Execute
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/go/bin/authx"]
