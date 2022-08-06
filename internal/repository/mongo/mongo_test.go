package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/danielMensah/user-management/internal/api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	createdAt = time.Time{}
	updatedAt = time.Time{}

	hexID1   = primitive.NewObjectIDFromTimestamp(createdAt).Hex()
	hexID2   = primitive.NewObjectIDFromTimestamp(updatedAt).Hex()
	nonHexID = primitive.NewObjectIDFromTimestamp(updatedAt).String()
)

func TestMain(m *testing.M) {
	var err error

	createdAt, err = time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		logrus.WithError(err).Fatal("failed to parse createdAt")
	}

	updatedAt, err = time.Parse(time.RFC3339, "2020-01-02T00:00:00Z")
	if err != nil {
		logrus.WithError(err).Fatal("failed to parse updatedAt")
	}

	logrus.Info("dates initialized")
	m.Run()
	os.Exit(0)
}

func TestNew(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	client := New(mt.DB)
	assert.NotNil(t, client)
}

func TestClient_GetUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name          string
		mockResponses []bson.D
		params        api.GetUsersParams
		mockError     bool
		expected      *[]api.User
		expectedErr   string
	}{
		{
			name: "successfully gets single user",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			mockResponses: []bson.D{
				{
					{"_id", hexID1},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				},
			},
			expected: &[]api.User{
				{
					Id:        hexID1,
					FirstName: "john",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@jd@mensah.com.com",
					Country:   "UK",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
		},
		{
			name: "successfully gets multiple users",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			mockResponses: []bson.D{
				{
					{"_id", hexID1},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				},
				{
					{"_id", hexID2},
					{"first_name", "jane"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				},
			},
			expected: &[]api.User{
				{
					Id:        hexID1,
					FirstName: "john",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@jd@mensah.com.com",
					Country:   "UK",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				{
					Id:        hexID2,
					FirstName: "jane",
					LastName:  "doe",
					Nickname:  "jd",
					Email:     "jd@jd@mensah.com.com",
					Country:   "UK",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
		},
		{
			name: "error retrieving users",
			mockResponses: []bson.D{
				{
					{"ok", "0"},
				},
			},
			mockError:   true,
			params:      api.GetUsersParams{},
			expected:    nil,
			expectedErr: errRetrieveFailed,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			if tt.mockError {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    0,
					Message: "error",
					Name:    "foo.bar.error",
				}))
			} else {
				r := make([]bson.D, 0)
				for i, b := range tt.mockResponses {
					id := mtest.FirstBatch
					if i > 0 {
						id = mtest.NextBatch
					}

					r = append(r, mtest.CreateCursorResponse(1, "foo.bar", id, b))
				}

				r = append(r, mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch))
				mt.AddMockResponses(r...)
			}

			c := &Client{
				db: mt.DB,
			}

			got, err := c.GetUsers(context.Background(), tt.params)

			if tt.expectedErr != "" {
				assert.Nil(t, got)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			teardown(mt)
		})
	}
}

func TestClient_CreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name         string
		mockResponse bson.D
		user         *api.UserCreateData
		expectedErr  string
	}{
		{
			name:         "can create user",
			mockResponse: bson.D{{"ok", 1}},
			user: &api.UserCreateData{
				FirstName: "john",
				LastName:  "doe",
				Nickname:  "jd",
				Email:     "jd@jd@mensah.com.com",
				Password:  "password",
				Country:   "UK",
			},
		},
		{
			name:         "error creating user",
			mockResponse: bson.D{{"ok", 0}},
			user: &api.UserCreateData{
				FirstName: "john",
				LastName:  "doe",
				Nickname:  "jd",
				Email:     "jd@jd@mensah.com.com",
				Password:  "password",
				Country:   "UK",
			},
			expectedErr: errInsertFailed,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			c := &Client{
				db: mt.DB,
			}

			mt.AddMockResponses(tt.mockResponse)

			got, err := c.CreateUser(context.Background(), tt.user)

			if tt.expectedErr != "" {
				assert.Empty(t, got)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestClient_UpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name               string
		id                 string
		data               *api.UserUpdateData
		mockUpdateResponse bson.D
		expected           *api.User
		expectedErr        string
	}{
		{
			name: "can update user",
			id:   hexID1,
			data: &api.UserUpdateData{
				FirstName: pstring("john"),
			},
			mockUpdateResponse: bson.D{
				{"ok", 1},
				{"value", bson.D{
					{"_id", hexID1},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				}},
			},
			expected: &api.User{
				Id:        hexID1,
				FirstName: "john",
				LastName:  "doe",
				Nickname:  "jd",
				Email:     "jd@jd@mensah.com.com",
				Country:   "UK",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		{
			name: "cannot update user",
			id:   hexID1,
			data: &api.UserUpdateData{
				FirstName: pstring("john"),
			},
			mockUpdateResponse: bson.D{},
			expectedErr:        errUpdateFailed,
		},
		{
			name:               "invalid id",
			id:                 nonHexID,
			data:               nil,
			mockUpdateResponse: bson.D{},
			expectedErr:        errConvertToObjectID,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			mt.AddMockResponses(tt.mockUpdateResponse)

			c := &Client{
				db: mt.DB,
			}

			got, err := c.UpdateUser(context.Background(), tt.id, tt.data)

			if tt.expectedErr != "" {
				assert.Nil(t, got)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name         string
		id           string
		mockResponse bson.D
		expectedErr  string
	}{
		{
			name: "can delete user",
			id:   hexID1,
			mockResponse: bson.D{
				{"ok", 1},
				{"value", bson.D{
					{"_id", hexID1},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"password", "password"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				}},
			},
			expectedErr: "",
		},
		{
			name:        "cannot delete user",
			id:          hexID1,
			expectedErr: errDeleteFailed,
		},
		{
			name:        "invalid id",
			id:          nonHexID,
			expectedErr: errConvertToObjectID,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			mt.AddMockResponses(tt.mockResponse)

			c := &Client{
				db: mt.DB,
			}

			err := c.DeleteUser(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
