# Authx

An implementation of auth service in Go. It's still a work-in-progress.

## Setup

### Git Hook

Add git hook to run a test and lint on every commit.

```bash
$ git config core.hooksPath .githooks
```

### Database

The auth service depends on Mongo for persistence.

* To run Mongo

  ```bash
  $ docker-compose up mongo
  $ docker exec -it mongo mongo -u nobody -p secrets --authenticationDatabase authx authx
  > show collections   # Show all collections in database authx
  ```

* Stop and remove the docker containers when done.

  ```bash
  $ docker-compose down
  ```

### Build, Test, and Run

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

* We need to install the `golangci-lint` before running the linter. Here's the standard installation:

  ```bash
  $ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1
  ```

  If you are on the Mac and has Homebrew installed, run:

  ```bash
  $ brew install golangci-lint
  ```

  Once you have `golangci-lint` installed, just run this command.

  ```bash
  $ make lint
  ```

* To format the Go code properly.

  ```bash
  $ make format  # Runs gofmt and rewrite the source code.
  ```

## Client

Use `curl` to test against the service.

```bash
curl --location --request POST 'http://localhost:8080/v1/signin' \
--header 'Content-Type: application/json' \
--data-raw '{
	"username": "chan",
	"password": "mypassword"
}'
```

## Testing

In the past, it really doesn't make sense to run database along with unit tests. The setup was slow and brittle. We mock the database in order to test code associated with the database.

With the advent of Docker, we get to set up the exact version of the database as a container on the development machine for unit testing. Even better, the setup is fast and flexible. The unit and integration tests in this project rely an actual database for testing.

### Setup Scripts and Makefile

Since we are including the database (running in containers) in both our unit and integration tests, we use [scripts/start-db-container.sh](scripts/start-db-container.sh) to orchestrate the following operations prior to running the tests that depend on the databse:

1. Spin up the database containers if they are not running. Skip start if the containers are running.
1. Wait till the database containers are ready for connection. After starting the container, we still need to wait for the TCP port to become open so that the unit or integration tests can start.

* `make test` - Run unit tests.
* `make end-db-container` - Tear down the database container.

## OAuth2

Here is [an in-depth description of OAuth2](docs).

## Troubleshooting

1. **Mongo docker container emits error ` no space left on device`.**

   This issue can be resolved by cleaning the old volumes with the following command:

   ```bash
   $ docker volume rm $(docker volume ls -qf dangling=true)
   ```

# Credits

* Gopher icon (used as favicon) by Renee French, CC BY 3.0 <https://creativecommons.org/licenses/by/3.0>, via Wikimedia Commons

# Reference

* [Github: MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
* [GoDocs: MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.4)
* [GoDocs: OAuth2 for Go](https://pkg.go.dev/golang.org/x/oauth2)
