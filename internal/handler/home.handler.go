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

	//get user from context
	const cKey = ContextKey("user")
	user := r.Context().Value(cKey)

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
	if user != nil {
		posts.User = user.(ClaimsCred)
		posts.ID = strconv.Itoa(int(posts.User.ID))
		posts.LoggedInUserId = posts.ID
	}

	posts.Posts = _posts
	common.ExecTemplate(w, "index.html", posts)

}

func Wrap(ID string, Title string) map[string]interface{} {
	return map[string]interface{}{
		"ID":    ID,
		"Title": Title,
	}
}
