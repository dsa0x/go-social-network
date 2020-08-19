package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dsa0x/go-social-network/common"
	"github.com/dsa0x/go-social-network/internal/model"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	post := model.Post{}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	post.CreatedBy = user.ID
	post.Author = user.Name

	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	ID, err := model.CreatePost(post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal server error"}`))
		return
	}

	resp := map[string]interface{}{"message": "Post created", "postId": ID}
	jsonPost, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonPost))

}

func CreatePost2(w http.ResponseWriter, r *http.Request, redirectPath string, deletePostID string) {
	const cKey = ContextKey("user")
	user := r.Context().Value(cKey)
	if user == nil {
		common.ExecTemplate(w, "index.html", user)
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

	if user != nil && r.Method == http.MethodPost && deletePostID == "" {
		user := user.(ClaimsCred)
		post := model.Post{}
		chip := r.FormValue("chip")
		post.Content = strings.TrimSpace(chip)
		if post.Content == "" {
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)

		}
		post.Author = user.Name
		post.CreatedBy = user.ID
		ID, err := model.CreatePost(post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}` + fmt.Sprint(ID)))
			return
		}
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)

	}
}

func DeletePost2(w http.ResponseWriter, r *http.Request, redirectPath string, deletePostID string) {
	//delete post from homepage
	const cKey = ContextKey("user")
	user := r.Context().Value(cKey)
	if user == nil {
		common.ExecTemplate(w, "index.html", user)
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
	if user != nil && r.Method == http.MethodPost && deletePostID != "" {
		user := user.(ClaimsCred)
		deletePostIDInt, _ := strconv.Atoi(deletePostID)
		_, err := model.DeletePost(uint(deletePostIDInt), user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
			return
		}
		// common.ExecTemplate(w, "index.html", posts)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		// return
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var postId, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad Request"}`))
		return
	}

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	ID, err := model.DeletePost(uint(postId), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Post not found"}`))
		return
	}

	resp := map[string]interface{}{"message": "Post deleted", "postId": ID}
	jsonPost, _ := json.Marshal(resp)
	// w.WriteHeader(http.StatusFound)
	w.Write([]byte(jsonPost))

}

func FetchAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	posts := []model.Post{}

	err := model.FetchPosts(&posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "An error occurred"}`))
		return
	}

	jsonPosts, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(jsonPosts))

}
func FetchUserPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	posts := []model.Post{}

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	err := model.FetchUserPosts(&posts, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "An error occurred"}`))
		return
	}

	jsonPosts, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(jsonPosts))

}
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	post := model.Post{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad Request"}`))
		return
	}

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey).(ClaimsCred)

	post, err = model.UpdatePost(post, uint(postID), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "An error occurred"}`))
		return
	}

	jsonPost, _ := json.Marshal(post)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(jsonPost))

}
