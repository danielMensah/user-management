package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestClient_GetUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	createdAt, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	updatedAt, err := time.Parse(time.RFC3339, "2020-01-02T00:00:00Z")
	require.NoError(t, err)

	tests := []struct {
		name        string
		batch       []bson.D
		params      api.GetUsersParams
		want        *[]api.User
		expectedErr string
	}{
		{
			name: "successfully gets single user",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			batch: []bson.D{
				{
					{"id", "abc"},
					{"firstName", "john"},
					{"lastName", "doe"},
					{"nickname", "jd"},
					{"email", "jd@challenge.com"},
					{"password", "password"},
					{"country", "UK"},
					{"createdAt", createdAt},
					{"updatedAt", updatedAt},
				},
			},
			want: &[]api.User{
				{
					Id:        "abc",
					FirstName: "john",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@challenge.com",
					Password:  "password",
					Country:   "UK",
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				},
			},
		},
		{
			name: "successfully gets multiple users",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			batch: []bson.D{
				{
					{"id", "abc"},
					{"firstName", "john"},
					{"lastName", "doe"},
					{"nickname", "jd"},
					{"email", "jd@challenge.com"},
					{"password", "password"},
					{"country", "UK"},
					{"createdAt", createdAt},
					{"updatedAt", updatedAt},
				},
				{
					{"id", "abc"},
					{"firstName", "jane"},
					{"lastName", "doe"},
					{"nickname", "jd"},
					{"email", "jd@challenge.com"},
					{"password", "password"},
					{"country", "UK"},
					{"createdAt", createdAt},
					{"updatedAt", updatedAt},
				},
			},
			want: &[]api.User{
				{
					Id:        "abc",
					FirstName: "john",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@challenge.com",
					Password:  "password",
					Country:   "UK",
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				},
				{
					Id:        "abc",
					FirstName: "jane",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@challenge.com",
					Password:  "password",
					Country:   "UK",
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				},
			},
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			r := make([]bson.D, 0)
			for i, b := range tt.batch {
				id := mtest.FirstBatch
				if i > 0 {
					id = mtest.NextBatch
				}

				r = append(r, mtest.CreateCursorResponse(1, "foo.bar", id, b))
			}

			r = append(r, mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch))
			mt.AddMockResponses(r...)

			c := &Client{
				db: mt.DB,
			}

			got, err := c.GetUsers(context.Background(), tt.params)

			if tt.expectedErr != "" {
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
