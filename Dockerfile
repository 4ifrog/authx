#####################
# --- Build Stage ---
#####################

# Pin version
FROM golang:1.15.6-alpine3.13 as build

# Install system dependencies for the build
RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    ncurses

# Support pretty print
ENV TERM=xterm-256color

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

FROM alpine:3.13

# Install tini - need for production for graceful shutdowns
RUN apk add --no-cache tini

# Install system runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    bash

# Copy binaries and config file over
WORKDIR /go/bin
USER nobody:nobody
COPY --from=build /go/src/bin/authx ./
COPY --from=build /go/src/bin/config.yaml ./
COPY --from=build /go/src/bin/static ./static
COPY --from=build /go/src/bin/templates ./templates

# Execute
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./authx"]
