# Authx

An implementation of an auth service written in Go.

## Setup

### Database

The auth service depends on Redis or Mongo for persistence.

* To run Redis:

  ```bash
  $ docker-compose up redis
  $ docker exec -it redis redis-cli
  127.0.0.1:6379> SELECT 0  # Use database 0
  127.0.0.1:6379> KEYS *    # Get all keys
  ```

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

* To run the linter.

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

# Reference

* [Github: MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
* [GoDocs: MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.4)

