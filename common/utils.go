package common

import (
	"encoding/json"
	"log"
	"net/http"

	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dsa0x/go-social-network/internal/config"
	"github.com/go-playground/validator/v10"
)

// Err error logger
func Err(err interface{}, message ...string) {
	if err != nil {
		log.Fatal(message[0], err)
	}
}

//hashPassword to hash password
func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate v
func Validate(s interface{}) map[string][]string {

	if err := validate.Struct(s); err != nil {

		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		//Validation errors occurred
		errors := make(map[string][]string)

		reflected := reflect.ValueOf(s)
		for _, err := range err.(validator.ValidationErrors) {

			field, _ := reflected.Type().FieldByName(err.StructField())
			var name string

			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				errors[name] = append(errors[name], "The "+name+" field is required")
				break
			case "email":
				errors[name] = append(errors[name], "The "+name+" field is not valid email")
				break
			case "eqfield":
				errors[name] = append(errors[name], "The "+name+" should be equal to the "+err.Param())
				break
			default:
				errors[name] = append(errors[name], "The "+name+" is invalid")
				break
			}
		}

		return errors
	}

	return nil

}

// CheckPasswordHash check the hash password match
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

//HashPassword to hash password
func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

// Jsonify js
func Jsonify(message interface{}) interface{} {
	// res := make(map[string]interface{})
	res, _ := json.Marshal(map[string]interface{}{"message": message})
	return res

}

// ExecTemplate execute template
func ExecTemplate(w http.ResponseWriter, template string, data interface{}) {

	err := config.Tpl.ExecuteTemplate(w, template, data)
	if err != nil {
		log.Fatalln("template didn't execute: ", template, err)
	}
}
