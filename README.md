# Authx

An implementation of an auth service written in Go.

## Setup

The auth service depends on Redis.

* To run the application as Docker containers

   ```bash
   $ docker-compose up  # Run both redis and authx
   ```

* To run the application directly.

   ```bash
   $ docker-compose up redis
   $ make run
   ```

* To test the application.

   ```bash
   $ make test
   ```

* To run the linter.

   ```bash
   $ make lint
   ```

* To format the Go code properly.

   ```bash
   $ make format  # Runs gofmt and rewrite the source code.
   ```
