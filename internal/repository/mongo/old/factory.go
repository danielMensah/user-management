package old

import (
	"context"
	"fmt"

	"github.com/danielMensah/faceit-challenge/internal/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	mongo mongoClient
	db    databaseClient
}

// New creates a new mongo repository
func New(connectionString, databaseName string) (repository.UserRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf(errConnectionFailed, err)
	}

	err = connect(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf(errConnectionFailed, err)
	}

	db := client.Database(databaseName)
	return &Client{
		mongo: client,
		db: repositoryClient{
			db: db,
		},
	}, nil
}

func connect(ctx context.Context, mongo mongoClient) error {
	err := mongo.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongo connection: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("mongo refused to connect: %v %w", ctx.Err(), err)
		default:
			err := mongo.Ping(ctx, readpref.Primary())
			if err == nil {
				logrus.Info("mongo is now connected")
				return nil
			}
		}
	}
}

func (c *Client) Close(ctx context.Context) error {
	return c.mongo.Disconnect(ctx)
}
