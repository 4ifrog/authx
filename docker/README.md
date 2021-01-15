# Docker

This directory contains docker-compose manifest files to run authx and one of its dependencies (mongo or redis).

## Authx and Mongo

Currently, Mongo is the preferred persistence engine for authx. To run this setup:

```bash
$ # This loads the docker-compose.yaml file.
$ docker-compose up
$ docker-compose down
```

If `docker-compose` is caching an older version of authx on your local machine, run the following:

```bash
$ docker-comoise up --build
```

## Authx and Redis

Currently, Redis is supported (but not sure for how much longer). To enable Redis in authx, we need to first modify the code to replace `StorageMongo` with `StorageRedis`. Finally run the following:

```bash
$ docker-compose -f docker-compose-redis.yaml up --build
$ docker-compose -f docker-compose-redis.yaml down
```
