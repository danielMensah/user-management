package old

import (
	"context"
	"fmt"
	"time"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	collectionUsers = "users"

	errConnectionFailed = "failed to connect to mongo: %w"
)

// repositoryClient represents the user repository contract
type repositoryClient struct {
	db mongoDatabase
}

type databaseClient interface {
	Collection(name string) mongoCollection
}

type mongoClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}

// mongoDatabase represents the mongo service contract. It is used for testing purposes.
type mongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

// Collection is a wrapper to get the mongo collection
func (r repositoryClient) Collection(name string) mongoCollection {
	return r.db.Collection(name)
}

// GetUsers returns a list of users
func (c *Client) GetUsers(ctx context.Context, params api.GetUsersParams) (*[]api.User, error) {
	opts := options.Find()
	opts.SetLimit(params.Limit)
	opts.SetSkip(params.Page)

	filter := bson.M{}

	collection := c.db.Collection(collectionUsers)
	cursor, err := collection.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users *[]api.User
	if err = cursor.All(ctx, users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Client) CreateUser(ctx context.Context, user *api.User) (string, error) {
	u := convertToMongoEntity(*user)

	result, err := c.db.Collection(collectionUsers).InsertOne(ctx, u)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	} else {
		return "", fmt.Errorf("unable to convert inserted id to object id")
	}
}

func (c *Client) UpdateUser(ctx context.Context, id string, data *api.UserUpdateData) (*api.User, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	panic("implement me")
}

func convertToMongoEntity(user api.User) UserCollection {
	return UserCollection{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Country:   user.Country,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
