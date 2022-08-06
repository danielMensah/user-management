package mongo

import (
	"context"
	"fmt"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"github.com/danielMensah/faceit-challenge/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionUsers = "users"
)

type Client struct {
	db *mongo.Database
}

func New(connectionString, databaseName string) (repository.UserRepository, error) {
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	db := conn.Database(databaseName)
	return &Client{db}, nil
}

func (c *Client) GetUsers(ctx context.Context, params api.GetUsersParams) (*[]api.User, error) {
	opts := options.Find()
	opts.SetLimit(params.Limit)
	opts.SetSkip(params.Page)

	filter := bson.M{}

	collection := c.db.Collection(collectionUsers)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]api.User, 1)
	if errC := cursor.All(ctx, &users); errC != nil {
		return nil, errC
	}

	return &users, nil
}

func (c *Client) CreateUser(ctx context.Context, user *api.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) UpdateUser(ctx context.Context, id string, data *api.UserUpdateData) (*api.User, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Client().Disconnect(ctx)
}
