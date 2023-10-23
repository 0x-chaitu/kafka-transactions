package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodb struct {
	DatabaseName string
	*mongo.Client
}

var (
	newClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
		return mongo.NewClient(opts...)
	}
	connect = func(ctx context.Context, client *mongo.Client) error {
		return client.Connect(ctx)
	}
	ping = func(ctx context.Context, client *mongo.Client) error {
		return client.Ping(ctx, nil)
	}
)

func Connect(ctx context.Context, host, database string, port int) (*Mongodb, error) {
	client, err := newClient(options.Client().ApplyURI(
		uri(host, port),
	))
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	err = connect(ctx, client)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	err = ping(ctx, client)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return &Mongodb{
		database,
		client,
	}, nil
}

func uri(host string, port int) string {
	const format = "mongodb://%s:%d"
	return fmt.Sprintf(format, host, port)
}
