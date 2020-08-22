package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dsa0x/go-social-network/common"
	"github.com/dsa0x/go-social-network/internal/model"

	"github.com/gorilla/mux"
)

type Profile struct {
	Username         string
	ID               string
	Following        int
	Followers        int
	IsFollower       bool
	Avatar           string
	MyProfile        bool
	Posts            []model.Post
	User             ClaimsCred
	LoggedInUserId   string
	LoggedInUsername string
	PostCount        int
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	const cKey = ContextKey("user")
	loggedUser := r.Context().Value(cKey)

	if loggedUser == nil {
		common.ExecTemplate(w, "index.html", loggedUser)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		// return
	}
	params := mux.Vars(r)
	redirectPath := "/user/" + params["id"]
	deletePostID := r.FormValue("postId")
	if loggedUser != nil && r.Method == http.MethodPost && deletePostID == "" {
		CreatePost2(w, r, redirectPath, deletePostID)
	}
	if loggedUser != nil && r.Method == http.MethodPost && deletePostID != "" {
		DeletePost2(w, r, "/", deletePostID)
	}

	profile := Profile{}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad request"}`))
		return
	}
	user, err := model.FindByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}
	loggedInUser := loggedUser.(ClaimsCred)
	if loggedInUser.ID == user.ID {
		profile.MyProfile = true
	}

	followings, err := model.CountFollowings(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}
	followers, err := model.CountFollowers(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}
	isFollower, _ := model.IsFollower(loggedInUser.ID, user.ID)

	posts := []model.Post{}
	err = model.FetchUserPosts(&posts, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "An error occurred"}`))
		return
	}
	count, err := model.CountUserPosts(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "An error occurred"}`))
		return
	}

	profile.Following = followings
	profile.Followers = followers
	profile.IsFollower = isFollower
	profile.ID = strconv.Itoa(int(user.ID))
	profile.Username = strings.Title(user.UserName)
	profile.Avatar = "profile_pic.svg"
	profile.Avatar = "/public/img/" + profile.Avatar
	profile.Posts = posts
	profile.User.Name = profile.Username
	profile.User.ID = loggedInUser.ID
	profile.LoggedInUserId = strconv.Itoa(int(loggedInUser.ID))
	profile.LoggedInUsername = loggedInUser.Name
	profile.PostCount = count

	// w.WriteHeader(http.StatusFound)
	common.ExecTemplate(w, "profile.html", profile)

}
func FollowUser(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad request"}`))
		return
	}
	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	_, err = model.Follow(user.ID, uint(id))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}
	http.Redirect(w, r, "/user/"+params["id"], http.StatusSeeOther)
	return
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad request"}`))
		return
	}
	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	_, err = model.Unfollow(user.ID, uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	http.Redirect(w, r, "/user/"+params["id"], http.StatusSeeOther)
	return
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	err := model.FetchUsers(&users)
	if err != nil {
		common.ExecTemplate(w, "users.html", nil)
		return
	}

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)
	userID := strconv.Itoa(int(user.ID))
	data := struct {
		LoggedInUserId string
		Users          []model.User
	}{
		userID,
		users,
	}

	common.ExecTemplate(w, "users.html", data)
}
