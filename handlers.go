package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

type UserInfoResponse struct {
	Id       int 		`json:"id"`
	UserInfo UserInfo 	`json:"userInfo"`
	Posts 	[]UserPost	`json:"posts"`
}
type UserInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserPost struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ServerErrorResponse is the error server response. It returns the reason, http status code,
// and request url for the response.
type ServerErrorResponse struct {
	StatusCode int 		`json:"statusCode"`
	RequestUrl string 	`json:"requestUrl"`
	Msg	string 			`json:"msg"`
}
func (in ServerErrorResponse) Error() string {
	return in.Msg
}
func NewServerErrorResponse(err error, requestUrl string, statusCode int) ServerErrorResponse {
	return ServerErrorResponse {
		StatusCode: statusCode,
		RequestUrl: requestUrl,
		Msg: err.Error(),
	}
}

type Handler struct {
	userClient user.Client
	logger zerolog.Logger
}

func NewHandler(client user.Client, logger zerolog.Logger) Handler {
	return Handler{
		userClient: client,
		logger: logger,
	}
}

// GetUserPostsHandler receives a userId and calls the UserAPI to fetch a user info along with the
// user's posts.
func(h Handler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	u, err := h.userClient.GetUserInfo(r.Context(), userID)
	if err != nil {
		h.handleErrorResponse(err, w, r)
		return
	}

	posts, err := h.userClient.GetUserPosts(r.Context(), userID)
	if err != nil {
		h.handleErrorResponse(err, w, r)
		return
	}

	userInfo := toUserInfoResponse(u, posts)
	if err = json.NewEncoder(w).Encode(userInfo); err != nil {
		h.handleErrorResponse(err, w, r)
		return
	}
}

// MiddlewareLogger is a http interceptor and logs each request that comes in and determines the log level based on
// the http status code that will be returned by the server.
func (h Handler) MiddlewareLogger(next http.Handler) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			// Recover and record stack traces in case of a panic
			if err := recover(); err != nil {
				h.logger.Error().
					Interface("recover_info", err).
					Bytes("debug_stack", debug.Stack()).
					Msgf("server error url=%s method=%s", r.URL, r.Method)
				http.Error(wrappedWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			logEvent := h.logger.Info()
			httpStatus := wrappedWriter.Status()
			if httpStatus >= http.StatusInternalServerError {
				logEvent = h.logger.Error()
			}
			body, _ := ioutil.ReadAll(r.Body)
			logEvent.
				Str("url", r.URL.Path).
				Str("method", r.Method).
				Str("body", string(body)).
				Int("status", httpStatus).
				Msgf("incoming request for %s", r.URL.Path)
		}()
		next.ServeHTTP(wrappedWriter, r)
	}
	return http.HandlerFunc(handlerFunc)
}

// handleErrorResponse logs and returns the appropriate http response code and response for errors from
// the client API or from an actual internal server error.
func (h Handler) handleErrorResponse(err error, w http.ResponseWriter, r *http.Request) {
	var apiClientError user.APIClientError
	if ok := errors.As(err,&apiClientError); ok {
		if apiClientError.StatusCode >= http.StatusInternalServerError {
			h.logger.Error().
				Err(apiClientError).
				Int("status", apiClientError.StatusCode).
				Msg("client API returned a server error")
		}
		w.WriteHeader(apiClientError.StatusCode)
		if err = json.NewEncoder(w).Encode(apiClientError); err != nil {
			h.logger.Error().Err(err).Msg("unable to encode APIClientError to JSON")
			http.Error(w, "unable to encode APIClientError to JSON", http.StatusInternalServerError)
			return
		}
		return
	}
	serverErrorResp := NewServerErrorResponse(err, r.URL.String(), http.StatusInternalServerError)
	h.logger.Error().Err(err).Msg("internal server error")
	w.WriteHeader(serverErrorResp.StatusCode)
	if err = json.NewEncoder(w).Encode(&serverErrorResp); err != nil {
		h.logger.Error().Err(err).Msg("unable to encode ServerErrorResponse to JSON")
		http.Error(w, "unable to encode ServerErrorResponse to JSON", http.StatusInternalServerError)
		return
	}
}

// toUserInfoResponse combines all the user.Post into a user.User
func toUserInfoResponse(user user.User, posts []user.Post) UserInfoResponse {
	userInfoResp := UserInfoResponse{
		Id: user.Id,
		UserInfo: UserInfo {
			Name: user.Name,
			Username: user.Username,
			Email:  user.Email,
		},
	}
	userPosts := make([]UserPost, 0, len(posts))
	for _, p := range posts {
		post := UserPost {
			Id: p.Id,
			Title: p.Title,
			Body: p.Body,
		}
		userPosts = append(userPosts, post)
	}
	userInfoResp.Posts = userPosts
	return userInfoResp
}
