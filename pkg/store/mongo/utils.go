package mongo

import (
	"context"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexOptions struct {
	isTTL bool
}

const (
	indexTimeout = 10 * time.Second
)

// maskDSN masks the sensitive username and password in the mongo dsn
func maskDSN(dsn string) string {
	re := regexp.MustCompile(`//.*@`)
	return re.ReplaceAllString(dsn, "//*****:*****@")
}

func containsKey(m primitive.M, key string) bool {
	_, ok := m[key]
	return ok
}

// containsChildKey returns true if param key matches the child key in { "key": { key: "someValue" }}
func containsChildKey(m primitive.M, key string) bool {
	val := m["key"]
	if val != nil {
		childMap := val.(primitive.M)
		_, ok := childMap[key]
		return ok
	}

	return false
}

// indexExists checks to see if index for `field` in `collection` already exist.
func indexExists(parent context.Context, collection *mongo.Collection, field string, opts ...*IndexOptions) (bool, error) {
	var isTTL bool
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		isTTL = opt.isTTL
	}

	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	var indexes []bson.M
	if err := cursor.All(context.Background(), &indexes); err != nil {
		return false, err
	}

	for _, m := range indexes {
		if isTTL {
			if containsKey(m, "expireAfterSeconds") && containsChildKey(m, field) {
				return true, nil
			}
		} else {
			if containsChildKey(m, field) {
				return true, nil
			}
		}
	}

	return false, nil
}

// createIndex creates an index associated with `field` in `collection`.
func createIndex(parent context.Context, collection *mongo.Collection, field string, opts ...*IndexOptions) (string, error) {
	var isTTL bool
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		isTTL = opt.isTTL
	}

	ok, err := indexExists(parent, collection, field, opts...)
	if err != nil {
		return "", err
	}
	if ok {
		// Skip index creation given the index already exist.
		return "", nil
	}

	// Create TTL index
	model := mongo.IndexModel{
		Keys: bson.D{
			{Key: field, Value: 1},
		},
	}

	if isTTL {
		// TODO: We may want to check that the field is of type Date as it is needed for TTL to work properly.
		model.Options = options.Index().SetExpireAfterSeconds(0)
	}

	ctx, cancel := context.WithTimeout(parent, atomicTimeout)
	defer cancel()

	ciOpts := options.CreateIndexes().SetMaxTime(indexTimeout)
	indexName, err := collection.Indexes().CreateOne(ctx, model, ciOpts)

	return indexName, err
}
