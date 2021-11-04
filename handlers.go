package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hooliganlin/simple-go-rest-api/user"
	"log"
	"net/http"
)

type UserInfoResponse struct {
	Id       int `json:"id"`
	UserInfo UserInfo `json:"userInfo"`
	Posts 	[]UserPost`json:"posts"`
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

type Handler struct {
	userClient user.Client
}

func NewHandler(client user.Client) Handler {
	return Handler{
		userClient: client,
	}
}

func(h Handler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	u, err := h.userClient.GetUserInfo(r.Context(), userID)
	if err != nil {
		//TODO return 500 since it's a legit error or we can wrap an error
		msg := fmt.Sprintf("Could not get user from uri=%s err=%v", r.RequestURI, err)
		log.Printf("[ERROR] %s", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	posts, err := h.userClient.GetUserPosts(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("Could not get posts for userId=$s uri=%s err=%v", r.RequestURI, err)
		log.Printf("[ERROR] %s", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	userInfo := toUserInfoResponse(u, posts)
	if err = json.NewEncoder(w).Encode(userInfo); err != nil {
		msg := fmt.Sprintf("Could not encode UserInfoResponse to JSON err=%v",  err)
		log.Printf("[ERROR] %s", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

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