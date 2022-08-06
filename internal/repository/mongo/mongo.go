package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/danielMensah/user-management/internal/api"
	"github.com/danielMensah/user-management/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionUsers = "users"

	errRetrieveFailed          = "failed to retrieve data from mongo"
	errCursorAllFailed         = "failed to use cursor to retrieve all data from mongo"
	errInsertFailed            = "failed to insert data into mongo"
	errConvertInsertedObjectID = "failed to convert inserted id to object id"
	errConvertToObjectID       = "failed to convert id string to object id"
	errUpdateFailed            = "failed to update user in mongo"
	errDeleteFailed            = "failed to delete user from mongo"
)

// Client represents a mongo client
type Client struct {
	db *mongo.Database
}

// New creates a new mongo repository client
func New(db *mongo.Database) repository.UserRepository {
	return &Client{db}
}

// GetUsers returns a list of users
func (c *Client) GetUsers(ctx context.Context, params api.GetUsersParams) (*[]api.User, error) {
	opts := options.Find()
	opts.SetLimit(params.Limit)
	opts.SetSkip(params.Page)
	opts.SetSort(bson.M{"createdAt": -1})

	filter := bson.M{}
	if params.Country != nil {
		filter["country"] = *params.Country
	}
	if params.Email != nil {
		filter["email"] = *params.Email
	}

	collection := c.db.Collection(collectionUsers)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errRetrieveFailed, err)
	}
	defer cursor.Close(ctx)

	users := make([]api.User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("%s: %w", errCursorAllFailed, err)
	}

	return &users, nil
}

// CreateUser creates a new user
func (c *Client) CreateUser(ctx context.Context, user *api.UserCreateData) (string, error) {
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	user.CreatedAt = &createdAt
	user.UpdatedAt = &updatedAt

	result, err := c.db.Collection(collectionUsers).InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errInsertFailed, err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return "", fmt.Errorf("%s: %s", errConvertInsertedObjectID, oid)
}

// UpdateUser updates a user
func (c *Client) UpdateUser(ctx context.Context, id string, data *api.UserUpdateData) (*api.User, error) {
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errConvertToObjectID, err)
	}

	updatedAt := time.Now().UTC()
	data.UpdatedAt = &updatedAt

	result := c.db.Collection(collectionUsers).FindOneAndUpdate(ctx, bson.M{"_id": pid}, bson.M{"$set": data})

	user := &api.User{}
	if err = result.Decode(user); err != nil {
		return nil, fmt.Errorf("%s with id '%s': %w", errUpdateFailed, id, err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", errConvertToObjectID, err)
	}

	result := c.db.Collection(collectionUsers).FindOneAndDelete(ctx, bson.M{"_id": pid})

	deletedUser := &api.User{}
	if err := result.Decode(deletedUser); err != nil {
		return fmt.Errorf("%s with id '%s': %w", errDeleteFailed, id, err)
	}

	return nil
}
