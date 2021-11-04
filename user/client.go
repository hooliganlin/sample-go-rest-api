package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company"`
}

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Client struct {
	baseURL string
	client  *http.Client
}

type Config struct {
	BaseURL string	`envconfig:"BASE_URL" default:"https://jsonplaceholder.typicode.com"`
}

func NewConfig() Config {
	var c Config
	err := envconfig.Process("userApi", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return c
}

func NewClient(c Config) Client {
	client := http.Client{
		Transport: http.DefaultTransport,
	}
	return Client{
		client:  &client,
		baseURL: c.BaseURL,
	}
}

func(c Client) GetUserInfo(ctx context.Context, userID string) (User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", c.baseURL, userID), nil)
	if err != nil {
		return User{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	//TODO handle different http statuses
	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (c Client) GetUserPosts(ctx context.Context, userID string) ([]Post, error){
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/posts?userId=%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	var posts []Post
	if err = json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}