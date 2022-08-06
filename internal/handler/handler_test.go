package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/danielMensah/user-management/internal/api"
	mongoRepo "github.com/danielMensah/user-management/internal/repository/mongo"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	createdAt = time.Time{}
	updatedAt = time.Time{}

	hexID = primitive.NewObjectID().Hex()
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

func setUpRequest(method, endpoint, body string) (echo.Context, *httptest.ResponseRecorder) {
	router := echo.New()
	request := httptest.NewRequest(method, endpoint, strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	echoContext := router.NewContext(request, response)
	return echoContext, response
}

func TestNewUserHandlers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	repo := mongoRepo.New(mt.DB)
	handlers := New(repo)
	assert.NotNil(t, handlers)
}

func TestService_GetUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name             string
		params           api.GetUsersParams
		mockResponses    []bson.D
		mockError        bool
		expectedStatus   int
		expectedResponse api.GetUsersResponse
		expectedErr      api.Error
	}{
		{
			name: "can get users",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			mockResponses: []bson.D{
				mtest.CreateSuccessResponse(bson.D{
					{"_id", hexID},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				}...),
			},
			expectedStatus: http.StatusOK,
			expectedResponse: api.GetUsersResponse{
				Users: &[]api.User{
					{
						Id:        hexID,
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
		},
		{
			name: "error getting users",
			params: api.GetUsersParams{
				Page:  0,
				Limit: 10,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: api.GetUsersResponse{},
			expectedErr: api.Error{
				Message: errGetUsers,
			},
			mockError: true,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

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

			repo := mongoRepo.New(mt.DB)
			h := &Handler{repo}

			ctx, response := setUpRequest(echo.GET, "/users", "")

			err := h.GetUsers(ctx, tt.params)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedErr.Message != "" {
				var responseBody api.Error
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedErr, responseBody)
			} else {
				var responseBody api.GetUsersResponse
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedResponse, responseBody)
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name           string
		body           string
		mockResponse   bson.D
		expectedStatus int
		expectedErr    api.Error
	}{
		{
			name:           "can create user",
			body:           `{"first_name":"john","last_name":"doe","nickname":"jd","email":"jd@jd@mensah.com.com","password":"password","country":"UK"}`,
			mockResponse:   bson.D{{"ok", 1}},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request",
			body:           `{"`,
			expectedStatus: http.StatusBadRequest,
			expectedErr: api.Error{
				Message: errParseBody,
			},
		},
		{
			name:           "error creating user",
			body:           `{"first_name":"john","last_name":"doe","nickname":"jd","email":"jd@jd@mensah.com.com","password":"password","country":"UK"}`,
			mockResponse:   bson.D{{"ok", 0}},
			expectedStatus: http.StatusInternalServerError,
			expectedErr: api.Error{
				Message: errCreateUser,
			},
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			mt.AddMockResponses(tt.mockResponse)

			repo := mongoRepo.New(mt.DB)
			s := &Handler{repo}

			ctx, response := setUpRequest(echo.POST, "/users", tt.body)

			err := s.CreateUser(ctx)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedErr.Message != "" {
				var responseBody api.Error
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedErr, responseBody)
			} else {
				var responseBody api.CreateUserResponse
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.NotEmpty(t, responseBody.Id)
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name             string
		id               string
		body             string
		mockResponse     bson.D
		expectedStatus   int
		expectedResponse api.User
		expectedErr      api.Error
	}{
		{
			name: "can update user",
			id:   hexID,
			body: `{"first_name":"john"}`,
			mockResponse: bson.D{
				{"ok", 1},
				{"value", bson.D{
					{"_id", hexID},
					{"first_name", "john"},
					{"last_name", "doe"},
					{"nickname", "jd"},
					{"email", "jd@jd@mensah.com.com"},
					{"country", "UK"},
					{"created_at", createdAt},
					{"updated_at", updatedAt},
				}},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: api.User{
				Id:        hexID,
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
			name:           "invalid request",
			id:             hexID,
			body:           `{"`,
			expectedStatus: http.StatusBadRequest,
			expectedErr: api.Error{
				Message: errParseBody,
			},
		},
		{
			name:           "error updating user",
			id:             hexID,
			body:           `{"first_name":"john"}`,
			mockResponse:   bson.D{{"ok", 0}},
			expectedStatus: http.StatusInternalServerError,
			expectedErr: api.Error{
				Message: errUpdateUser,
			},
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			mt.AddMockResponses(tt.mockResponse)

			repo := mongoRepo.New(mt.DB)
			s := &Handler{repo}

			ctx, response := setUpRequest(echo.PUT, "/users/:hexID", tt.body)

			err := s.UpdateUser(ctx, tt.id)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedErr.Message != "" {
				var responseBody api.Error
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedErr, responseBody)
			} else {
				var responseBody api.User
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedResponse, responseBody)
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name           string
		id             string
		mockResponse   bson.D
		expectedStatus int
		expectedErr    api.Error
	}{
		{
			name: "can delete user",
			id:   hexID,
			mockResponse: bson.D{
				{"ok", 1},
				{"value", bson.D{
					{"_id", hexID},
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
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "error deleting user",
			id:             hexID,
			mockResponse:   bson.D{{"ok", 0}},
			expectedStatus: http.StatusInternalServerError,
			expectedErr: api.Error{
				Message: errDeleteUser,
			},
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			defer teardown(mt)

			mt.AddMockResponses(tt.mockResponse)

			repo := mongoRepo.New(mt.DB)
			s := &Handler{repo}

			ctx, response := setUpRequest(echo.PUT, "/users/:hexID", "")

			err := s.DeleteUser(ctx, tt.id)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedErr.Message != "" {
				var responseBody api.Error
				err = json.Unmarshal(response.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedErr, responseBody)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
