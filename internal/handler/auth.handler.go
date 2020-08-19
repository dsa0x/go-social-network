package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dsa0x/go-social-network/common"
	"github.com/dsa0x/go-social-network/internal/config"
	"github.com/dsa0x/go-social-network/internal/model"
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

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var errs = make(Error)

	if r.Method == http.MethodPost {
		creds := Credentials{}
		email := r.FormValue("email")
		password := r.FormValue("password")
		creds.Email = email
		creds.Password = password

		if err := common.Validate(creds); err != nil {
			// message, _ := json.Marshal(err)
			w.WriteHeader(http.StatusBadRequest)
			// w.Write(message)
			// errs = err
			common.ExecTemplate(w, "login.html", err)
			return
		}

		user, err := model.FindOne(creds.Email)
		if err != nil {
			// message, _ := json.Marshal(err)
			w.WriteHeader(http.StatusBadRequest)
			errs["message"] = append(errs["message"], err.Error())
			common.ExecTemplate(w, "login.html", errs)
			// w.Write(message)
			return
		}

		matched := common.CheckPasswordHash(creds.Password, user.Password)
		if err != nil || !matched {
			log.Println("Username or password Incorrect: ")
			w.WriteHeader(http.StatusForbidden)
			// w.Write([]byte(`{"message": "Username or password Incorrect"}`))
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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte(`{"message": "An error occured"}`))
			errs["mismatch"] = append(errs["mismatch"], "An error occured. please try again")
			common.ExecTemplate(w, "login.html", errs)
			return
		}

		cookie := &http.Cookie{
			Name: "session", Value: tokenString,
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	common.ExecTemplate(w, "login.html", nil)

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
			// w.Write([]byte(`{"message": "Passwords do not match"}`))
			errs["mismatch"] = append(errs["mismatch"], "Passwords do not match")
			common.ExecTemplate(w, "signup.html", errs)
			return
		}
		id, err := model.CreateUser(user)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("User with id %d created", id)

		// redirect
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

	common.ExecTemplate(w, "signup.html", nil)

}

//Auth authenticate user
func Auth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		reqToken, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "Invalid Cookies"}`))
			return
		}
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(reqToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "Invalid Signature"}`))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid Authorization"}`))
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "Invalid Authorization"}`))
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
		common.ExecTemplate(w, "header.html", data)
		f(w, r)
	}

}
