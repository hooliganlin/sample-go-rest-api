package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserPostsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/user-posts/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	u := user.User {
		Id: 1,
		Name: "my first name",
		Username: "this_username",
		Email: "first@example.com",
	}
	posts := []user.Post{
		{
			UserId: 1,
			Id: 1,
			Title: "my first title",
			Body: "my first body",
		},
		{
			UserId: 1,
			Id: 2,
			Title: "my second title",
			Body: "my second body",
		},
	}

	mockContext := mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})

	t.Run("successful response", func(t *testing.T) {
		mockClient := new(MockUserClient)
		logger := zerolog.New(io.Discard)
		handler := NewHandler(mockClient, logger)

		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)

		mockClient.On("GetUserInfo", mockContext, "1").Return(u, nil)
		mockClient.On("GetUserPosts", mockContext, "1").Return(posts, nil)

		handler.GetUserPostsHandler(recorder, req)
		expectedResult := toUserInfoResponse(u, posts)
		var result UserInfoResponse
		if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedResult, result)
	})

	t.Run("getUserInfo failure", func(t *testing.T) {
		mockClient := new(MockUserClient)
		logger := zerolog.New(io.Discard)
		handler := NewHandler(mockClient, logger)

		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusNotFound)
		resp := recorder.Result()

		err := user.NewAPIClientError(resp, req)
		mockClient.On("GetUserInfo", mockContext, "1").Return(user.User{}, err)

		handler.GetUserPostsHandler(recorder, req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.JSONEq(t,
			`{"statusCode":404,"msg":"API returned an invalid or empty response","url":"/v1/user-posts/1"}`,
			recorder.Body.String())
	})

	t.Run("getUserPosts failure", func(t *testing.T) {
		mockClient := new(MockUserClient)
		logger := zerolog.New(io.Discard)
		handler := NewHandler(mockClient, logger)

		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusNotFound)
		resp := recorder.Result()

		err := user.NewAPIClientError(resp, req)
		mockClient.On("GetUserInfo", mockContext, "1").Return(u, nil)
		mockClient.On("GetUserPosts", mockContext, "1").Return([]user.Post{}, err)

		handler.GetUserPostsHandler(recorder, req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.JSONEq(t,
			`{"statusCode":404,"msg":"API returned an invalid or empty response","url":"/v1/user-posts/1"}`,
			recorder.Body.String())
	})
}

func TestToUserInfoResponse(t *testing.T) {
	u := user.User {
		Id:       1,
		Name:     "Bob Loblaw",
		Username: "bob",
		Email:    "bob@laywer.com",
		Phone:   "123-456-1234",
		Website: "www.bob.com",
	}
	posts := []user.Post{
		{
			UserId: 1,
			Id: 1,
			Title: "Lorem Ipsum",
			Body: "Brewing coffee",
		},
		{
			UserId: 1,
			Id: 2,
			Title: "The second title",
			Body: "Brewing coffee again",
		},
	}
	t.Run("success", func(t *testing.T) {
		result := toUserInfoResponse(u, posts)
		assert.Equal(t, UserInfoResponse{
			Id:       1,
			UserInfo: UserInfo{
				Name: u.Name,
				Username: u.Username,
				Email: u.Email,
			},
			Posts:    []UserPost{
				{
					Id: posts[0].Id,
					Title: posts[0].Title,
					Body: posts[0].Body,
				},
				{
					Id: posts[1].Id,
					Title: posts[1].Title,
					Body: posts[1].Body,
				},
			},
		}, result)
	})

	t.Run("no posts", func(t *testing.T) {
		result := toUserInfoResponse(u, nil)
		assert.Equal(t, UserInfoResponse{
			Id:       1,
			UserInfo: UserInfo{
				Name: u.Name,
				Username: u.Username,
				Email: u.Email,
			},
			Posts: []UserPost{},
		}, result)
	})
}

func TestHandleErrorResponse(t *testing.T) {
	mockClient := new(MockUserClient)
	out := &bytes.Buffer{}
	logger := zerolog.New(out)
	handler := NewHandler(mockClient, logger)

	req := httptest.NewRequest(http.MethodGet, "/v1/user-posts/1", nil)
	t.Run("ApiClientError", func(t *testing.T) {
		t.Run("http 500 status code", func(t *testing.T) {
			res := callErrorHandlerWithApiClientError(handler, http.StatusInternalServerError, req)

			logResult := convertJSONToMap(out)
			assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
			assert.JSONEq(t,
				`{"statusCode":500,"msg":"API returned an invalid or empty response","url":"/v1/user-posts/1"}`,
				res.Body.String())
			assert.EqualValues(t, http.StatusInternalServerError, logResult["status"])
			assert.Equal(t, "client API returned a server error", logResult["message"])
			out.Reset()
		})
		t.Run("http non 500 status code", func(t *testing.T) {
			res := callErrorHandlerWithApiClientError(handler, http.StatusNotFound, req)
			assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
			assert.JSONEq(t,
				`{"statusCode":404,"msg":"API returned an invalid or empty response","url":"/v1/user-posts/1"}`,
				res.Body.String())
			out.Reset()
		})
	})

	t.Run("serverErrorResponse", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusInternalServerError)
		res := recorder.Result()
		defer res.Body.Close()

		handler.handleErrorResponse(errors.New("Oh no no no"), recorder, req)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.JSONEq(t, `{"statusCode":500,"requestUrl":"/v1/user-posts/1","msg":"Oh no no no"}`, recorder.Body.String())
	})

}

func callErrorHandlerWithApiClientError(handler Handler, inStatus int, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(inStatus)
	res := recorder.Result()
	defer res.Body.Close()

	apiClientErr := user.NewAPIClientError(res, req)
	handler.handleErrorResponse(apiClientErr, recorder, req)
	return recorder
}

func convertJSONToMap(out *bytes.Buffer) map[string]interface{}{
	var result map[string]interface{}
	_ = json.Unmarshal(out.Bytes(), &result)
	return result
}

type MockUserClient struct {
	mock.Mock
}
func (m *MockUserClient) GetUserInfo(ctx context.Context, userID string) (user.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *MockUserClient) GetUserPosts(ctx context.Context, userID string) ([]user.Post, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]user.Post), args.Error(1)
}