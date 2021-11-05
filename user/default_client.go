package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Config struct {
	BaseURL string	`envconfig:"BASE_URL" default:"https://jsonaplaceholder.typicode.com"`
}

func NewConfig() Config {
	var c Config
	if err := envconfig.Process("userApi", &c); err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return c
}

type DefaultClient struct {
	baseURL string
	client  *http.Client
}

func NewDefaultClient(c Config) Client {
	client := http.Client{
		Transport: http.DefaultTransport,
	}
	return DefaultClient{
		client:  &client,
		baseURL: c.BaseURL,
	}
}

// GetUserInfo fetches user information from the User API
func(c DefaultClient) GetUserInfo(ctx context.Context, userID string) (User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", c.baseURL, userID), nil)
	if err != nil {
		return User{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return User{}, err
	}
	if err = checkResponse(resp); err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUserPosts fetches posts for a user from the UserPost API
func (c DefaultClient) GetUserPosts(ctx context.Context, userID string) ([]Post, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/posts?userId=%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err = checkResponse(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []Post
	if err = json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}
