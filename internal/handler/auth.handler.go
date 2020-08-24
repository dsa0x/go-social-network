package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dsa0x/go-social-network/common"
	"github.com/dsa0x/go-social-network/internal/config"
	"github.com/dsa0x/go-social-network/internal/model"
	"github.com/gorilla/securecookie"
)

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	Name  string `json:"username"`
	jwt.StandardClaims
}

type ClaimsCred struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	Name  string `json:"username"`
}

type ContextKey string

type Error map[string][]string

var jwtKey = []byte(config.Env.JwtKey)
var hashKey = []byte(config.Env.HashKey)
var secureCookie = securecookie.New(hashKey, nil)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var errs = make(Error)

	if r.Method == http.MethodPost {
		creds := Credentials{}
		email := r.FormValue("email")
		password := r.FormValue("password")
		creds.Email = strings.ToLower(email)
		creds.Password = password

		if err := common.Validate(creds); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err, "::Validation failed")
			errs["mismatch"] = append(errs["mismatch"], "Username or password incorrect")
			common.ExecTemplate(w, "login.html", err)
			return
		}

		user, err := model.FindOne(creds.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err, "::No user found")
			errs["mismatch"] = append(errs["mismatch"], "Username or password incorrect")
			common.ExecTemplate(w, "login.html", errs)
			return
		}

		matched := common.CheckPasswordHash(creds.Password, user.Password)
		if err != nil || !matched {
			log.Println(err, user.ID, user.UserName, "::password mismatch")
			w.WriteHeader(http.StatusForbidden)
			errs["mismatch"] = append(errs["mismatch"], "Username or password incorrect")
			common.ExecTemplate(w, "login.html", errs)
			return
		}

		expirationTime := time.Now().Add(30 * time.Minute)

		claims := &Claims{
			Email: user.Email,
			ID:    user.ID,
			Name:  user.UserName,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			log.Println(err, "token signature error")
			w.WriteHeader(http.StatusInternalServerError)
			errs["mismatch"] = append(errs["mismatch"], "An error occured. please try again")
			common.ExecTemplate(w, "login.html", errs)
			return
		}

		encoded, err := secureCookie.Encode("session", tokenString)
		expire := time.Now().Add(24 * time.Hour)
		cookie := &http.Cookie{
			Name: "session", Value: encoded, HttpOnly: true, Expires: expire,
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	common.ExecTemplate(w, "login.html", nil)

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var errs = make(Error)
	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")
		cpassword := r.FormValue("cpassword")
		username := r.FormValue("username")
		user := model.User{}
		user.Email = email
		user.Password = password
		user.ConfirmPassword = cpassword
		user.UserName = username

		if user.Password == "" || user.Password != user.ConfirmPassword {
			w.WriteHeader(http.StatusNotFound)
			errs["mismatch"] = append(errs["mismatch"], "Passwords do not match")
			common.ExecTemplate(w, "signup.html", errs)
			return
		}

		id, err := model.CreateUser(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			if strings.Contains(err.Error(), `unique constraint "users_user_name_key"`) {
				errs["mismatch"] = append(errs["mismatch"], "Username already exists")
			} else if strings.Contains(err.Error(), `unique constraint "users_email_key"`) {
				errs["mismatch"] = append(errs["mismatch"], "Email already exists")
			} else {
				errs["mismatch"] = append(errs["mismatch"], err.Error())
			}
			common.ExecTemplate(w, "signup.html", errs)
			return
		}

		log.Printf("User with id %d created", id)

		// redirect
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

	common.ExecTemplate(w, "signup.html", nil)

}

func Logout(w http.ResponseWriter, r *http.Request) {

	expire := time.Now().Add(-7 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name: "session", Value: "", Expires: expire,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

//Auth authenticate user
func Auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println(err, "Invalid Session")
			http.Redirect(w, r, "/guest", http.StatusSeeOther)
			return
		}
		var reqToken string
		if err = secureCookie.Decode("session", cookie.Value, &reqToken); err != nil {
			log.Println(err, "Unable to decode session")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			log.Println(err, "Invalid Authorization")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !token.Valid {
			log.Println(err, "Invalid Token")
			w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		const cKey = ContextKey("user")
		ctx := context.WithValue(r.Context(), cKey, ClaimsCred{Email: claims.Email, ID: claims.ID, Name: claims.Name})

		var data struct {
			ID    string
			Name  string
			Title string
		}
		data.ID = fmt.Sprint(claims.ID)
		data.Name = claims.Name

		r = r.WithContext(ctx)
		f(w, r)
	}

}
