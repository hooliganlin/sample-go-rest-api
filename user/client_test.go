package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	t.Run("http 200 response", func(t *testing.T) {
		expectedUser := User{
			Id: 1,
			Name: "Yolanda",
			Username: "thunder_chunky",
			Email: "yolanda@example.com",
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewEncoder(w).Encode(expectedUser); err != nil {
				t.Error(err, "could not encode user to JSON")
			}
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer testServer.Close()

		client := NewClient(Config{
			BaseURL: testServer.URL,
		})

		u, err := client.GetUserInfo(context.Background(), "user_1")
		if err != nil {
			t.Error(err, "could not call getUserInfo")
		}
		assert.Equal(t, expectedUser, u)
	})

	t.Run("non http 2xx response", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer testServer.Close()

		client := NewClient(Config{
			BaseURL: testServer.URL,
		})

		u, err := client.GetUserInfo(context.Background(), "user_1")
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("API response error statusCode=404 body= url=%s/users/user_1", testServer.URL), err.Error())
		assert.Equal(t, User{}, u)
	})
}

func TestGetUserPosts(t *testing.T) {
	t.Run("http 200 response", func(t *testing.T) {
		expectedPosts := []Post{
			{
				Id: 1,
				UserId: 1,
				Title: "Can do!",
				Body: "Lorem ipsum here and there",
			},
			{
				Id: 2,
				UserId: 1,
				Title: "Cannot do!",
				Body: "the other body",
			},
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewEncoder(w).Encode(expectedPosts); err != nil {
				t.Error(err, "could not encode posts to JSON")
			}
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer testServer.Close()

		client := NewClient(Config{
			BaseURL: testServer.URL,
		})

		posts, err := client.GetUserPosts(context.Background(), "user_1")
		if err != nil {
			t.Error(err, "could not call getUserInfo")
		}
		assert.Equal(t, expectedPosts, posts)
	})

	t.Run("non http 2xx response", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer testServer.Close()

		client := NewClient(Config{
			BaseURL: testServer.URL,
		})

		_, err := client.GetUserPosts(context.Background(), "user_1")
		assert.Error(t, err)
	})
}