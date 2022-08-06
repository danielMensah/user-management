package old

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserCollection represents a user in mongoDB
type UserCollection struct {
	Id        uuid.UUID `bson:"id,omitempty"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	Nickname  string    `bson:"nickname"`
	Country   string    `bson:"country"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

type mongoCollectionClient struct {
	col mongoCollection
}

type mongoCursor interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
}

type mongoCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

func (c mongoCollectionClient) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (mongoCursor, error) {
	return c.col.Find(ctx, filter, opts...)
}

func (c mongoCollectionClient) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.col.InsertOne(ctx, document, opts...)
}
