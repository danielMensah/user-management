package old

import (
	"context"
	"testing"
	"time"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockRepositoryClient struct {
	repositoryClient
	mock.Mock
}

// Collection is a mock of this method used for testing
func (m *mockRepositoryClient) Collection(name string) mongoCollection {
	args := m.Called(name)

	collection, ok := args.Get(0).(mongoCollection)
	if !ok {
		collection = nil
	}

	return collection
}

type mockMongoCollection struct {
	mongoCollection
	mock.Mock
}

// Find is a mock of this method used for testing
func (m *mockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (mongoCursor, error) {
	args := m.Called(ctx, filter, opts)
	cursor, ok := args.Get(0).(mongoCursor)
	if !ok {
		cursor = nil
	}
	return cursor, args.Error(1)
}

type mockMongoCursor struct {
	mongoCursor
	mock.Mock
}

func (m *mockMongoCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

func (m *mockMongoCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRepository_GetUsers(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()
	ctx := context.Background()

	tests := []struct {
		name               string
		params             api.GetUsersParams
		repo               *mockRepositoryClient
		collection         *mockMongoCollection
		cursor             *mockMongoCursor
		expectedMongoMocks func(t *testing.T, repo *mockRepositoryClient, col *mockMongoCollection, cur *mockMongoCursor)
		want               *[]api.User
		expectedErr        string
	}{
		{
			name: "can fetch users",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			repo:       &mockRepositoryClient{},
			collection: &mockMongoCollection{},
			cursor:     &mockMongoCursor{},
			want: &[]api.User{
				{
					Id:        "64194491-df5a-4ed2-9175-cf389fd2e640",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "JD",
					Email:     "jd@faceitchallenge.com",
					Password:  "password",
					Country:   "UK",
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				},
			},
			expectedMongoMocks: func(t *testing.T, repo *mockRepositoryClient, col *mockMongoCollection, cur *mockMongoCursor) {
				opts := options.Find()
				opts.SetSkip(0)
				opts.SetLimit(10)

				var users *[]api.User

				repo.On("Collection", collectionUsers).Return(col)
				col.On("Find", ctx, bson.M{}, []*options.FindOptions{(*options.FindOptions)(nil)}).Return(cur, nil)
				cur.On("All", ctx, users).Return(nil)
				cur.On("Close", ctx).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				db: tt.repo,
			}

			if tt.expectedMongoMocks != nil {
				tt.expectedMongoMocks(t, tt.repo, tt.collection, tt.cursor)
			}

			users, err := c.GetUsers(ctx, tt.params)

			if tt.expectedErr != "" {
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, users)
			}

			if tt.expectedMongoMocks != nil {
				tt.repo.AssertExpectations(t)
				tt.collection.AssertExpectations(t)
				tt.cursor.AssertExpectations(t)
			}
		})
	}
}
