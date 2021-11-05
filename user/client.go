package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
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
		Suite   string `json:"suite,omitempty"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo,omitempty"`
	} `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Website string `json:"website,omitempty"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase,omitempty"`
		Bs          string `json:"bs,omitempty"`
	} `json:"company,omitempty"`
}

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type DefaultClient struct {
	baseURL string
	client  *http.Client
}

type Client interface {
	GetUserInfo(ctx context.Context, userID string) (User, error)
	GetUserPosts(ctx context.Context, userID string) ([]Post, error)
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

// CheckResponse checks an API response and returns and error if
// API returns an error response.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	return NewAPIClientError(resp, resp.Request)
}

// APIClientError is an error container for any errors from the API.
type APIClientError struct {
	StatusCode int		`json:"statusCode"`
	Body string			`json:"body,omitempty"`
	Msg string			`json:"msg,omitempty"`
	URL string			`json:"url"`
}

func NewAPIClientError(resp *http.Response, req *http.Request) APIClientError {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIClientError{}
	}
	j := string(body); if j != "{}" && j != "" {
		return APIClientError{
			StatusCode: resp.StatusCode,
			Body: string(body),
			URL: req.URL.String(),
		}
	}
	return APIClientError{
		StatusCode: resp.StatusCode,
		Msg: "API returned an invalid or empty response",
		URL: req.URL.String(),
	}
}

func (in APIClientError) Error() string {
	return fmt.Sprintf("API response error statusCode=%d body=%s url=%s", in.StatusCode, in.Body, in.URL)
}