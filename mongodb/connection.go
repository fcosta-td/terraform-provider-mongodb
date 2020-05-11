package mongodb

import (
	"context"
	"log"
	"time"
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type CreateResult struct {
	Ok int `bson:"ok"`
}

func (c *CreateResult) IsOk() bool {
	return c.Ok == 1
}

type Client struct {
	client	*mongo.Client
	context *context.Context
}

// RedactMongoUri removes login and password from mongoUri.
func RedactMongoUri(uri string) string {
	if strings.HasPrefix(uri, "mongodb://") && strings.Contains(uri, "@") {
		if strings.Contains(uri, "ssl=true") {
			uri = strings.Replace(uri, "ssl=true", "", 1)
		}

		cStr, err := connstring.Parse(uri)
		if err != nil {
			log.Printf("[Err] Cannot parse mongodb server url: %s", err)
			return "unknown/error"
		}

		if cStr.Username != "" && cStr.Password != "" {
			uri = strings.Replace(uri, cStr.Username, "****", 1)
			uri = strings.Replace(uri, cStr.Password, "****", 1)
			return uri
		}
	}
	return uri
}

// Creates connection
func NewClient(connectionString string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)
	ops := options.Client().
		ApplyURI(connectionString).
		SetDirect(true).
		SetReadPreference(readpref.Primary()).
		SetAppName("terraform-provider-mongodb")

	c, err := mongo.Connect(ctx,ops)

	if err != nil {
		log.Printf("[ERR] Cannot connect to server using url %s: %s", RedactMongoUri(connectionString), err)
		return nil, err
	}

	client := &Client{
		client:  c,
		context: &ctx,
	}

	return client, nil
}

func CheckUserPassword(database string, username string, password string) (string) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)
	ops := options.Client().
		ApplyURI(mongo_uri).
		SetDirect(true).
		SetReadPreference(readpref.Secondary()).
		SetAppName("terraform-provider-mongodb").
		SetAuth(options.Credential{
			AuthSource: database, Username: username, Password: password,
		 })

	c, err := mongo.Connect(ctx,ops)

	if err != nil {
		log.Printf("Failed to connect to mongo uri")
	}

	err = c.Ping(ctx, readpref.Primary())
	c.Disconnect(ctx)

	if err != nil {
		return "change_me"
	} else {
		return password
	}
}
