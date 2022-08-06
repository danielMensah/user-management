package old

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoCollectionClient struct {
	mongoCollectionClient
	mock.Mock
}

// Find is a mock of this method used for testing
func (m *mockMongoCollectionClient) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	args := m.Called(ctx, filter, opts)
	result, ok := args.Get(0).(*mongo.Cursor)
	if !ok {
		result = nil
	}
	return result, args.Error(1)
}

func Test_mongoCollectionClient_Find(t *testing.T) {
	filter := bson.M{}
	ctx := context.Background()

	tests := []struct {
		name               string
		collection         *mockMongoCollectionClient
		expectedMongoMocks func(t *testing.T, col *mockMongoCollectionClient)
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:       "successfully executes a find operation",
			collection: &mockMongoCollectionClient{},
			expectedMongoMocks: func(t *testing.T, col *mockMongoCollectionClient) {
				col.On("Find", ctx, filter, mock.Anything).Return(nil, nil)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mongoCollectionClient{
				col: tt.collection,
			}

			if tt.expectedMongoMocks != nil {
				tt.expectedMongoMocks(t, tt.collection)
			}

			_, err := c.Find(ctx, filter, options.Find())
			assert.NoError(t, err)

			if tt.expectedMongoMocks != nil {
				tt.collection.AssertExpectations(t)
			}
		})
	}
}
