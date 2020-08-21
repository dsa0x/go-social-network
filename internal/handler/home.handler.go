package handler

import (
	"net/http"
	"strconv"

	"github.com/dsa0x/go-social-network/internal/model"

	"github.com/dsa0x/go-social-network/common"
)

type HomePosts struct {
	Posts          []model.Post
	User           ClaimsCred
	ID             string
	LoggedInUserId string
}

// Home function for home handler
func Home(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")

	const cKey = ContextKey("user")
	user := r.Context().Value(cKey)
	if user == nil {
		common.ExecTemplate(w, "index.html", user)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		// return
	}
	deletePostID := r.FormValue("postId")
	if user != nil && r.Method == http.MethodPost && deletePostID == "" {
		CreatePost2(w, r, "/", deletePostID)
	}
	if user != nil && r.Method == http.MethodPost && deletePostID != "" {
		DeletePost2(w, r, "/", deletePostID)
	}

	// fetch all posts for homepage
	_posts := []model.Post{}
	err := model.FetchPosts(&_posts)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	posts := HomePosts{}
	posts.User = user.(ClaimsCred)
	posts.ID = strconv.Itoa(int(posts.User.ID))
	// posts.ID = posts.User.ID//
	posts.LoggedInUserId = posts.ID
	posts.Posts = _posts
	common.ExecTemplate(w, "index.html", posts)

}

func Wrap(ID string, Title string) map[string]interface{} {
	return map[string]interface{}{
		"ID":    ID,
		"Title": Title,
	}
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
