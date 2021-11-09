package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hooliganlin/simple-go-rest-api/cache"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	userPostCacheKeyPrefix = "posts-user"
	userCacheKeyPrefix = "user"
)

type Config struct {
	BaseURL string	`envconfig:"BASE_URL" default:"https://jsonplaceholder.typicode.com"`
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
	cache 	cache.Cache
}

func NewDefaultClient(c Config, cache cache.Cache) Client {
	client := http.Client{
		Transport: cleanhttp.DefaultPooledTransport(),
	}
	return DefaultClient{
		client:  &client,
		baseURL: c.BaseURL,
		cache: cache,
	}
}

// GetUserInfo fetches user information from the User API
func(c DefaultClient) GetUserInfo(ctx context.Context, userID string) (User, error) {
	// check cache first
	cacheKey := userCacheKey(userID)
	if u, ok := c.cache.Get(cacheKey); ok {
		return u.(User), nil
	}

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
	c.cache.Set(cacheKey, User{})
	return user, nil
}

// GetUserPosts fetches posts for a user from the UserPost API
func (c DefaultClient) GetUserPosts(ctx context.Context, userID string) ([]Post, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/posts?userId=%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, err
	}
	cacheKey := userPostsCacheKey(userID)
	if p, ok := c.cache.Get(cacheKey); ok {
		return p.([]Post), nil
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
	c.cache.Set(cacheKey, posts)
	return posts, nil
}

func userCacheKey(userID string) string {
	return fmt.Sprintf("%s-%s", userCacheKeyPrefix, userID)
}
func userPostsCacheKey(userID string) string {
	return fmt.Sprintf("%s-%s", userPostCacheKeyPrefix, userID)
}