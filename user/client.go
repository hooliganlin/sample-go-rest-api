package user

import (
	"context"
	"fmt"
	"io/ioutil"
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

type Client interface {
	GetUserInfo(ctx context.Context, userID string) (User, error)
	GetUserPosts(ctx context.Context, userID string) ([]Post, error)
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
	if j := string(body); j != "{}" && j != "" {
		return APIClientError{
			StatusCode: resp.StatusCode,
			Body: j,
			Msg: "API returned an error",
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