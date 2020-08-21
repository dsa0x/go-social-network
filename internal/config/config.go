package config

import (
	"html/template"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type envVars struct {
	Dbhost     string `required:"true" envconfig:"DB_HOST"`
	Dbport     string `required:"true" envconfig:"DB_PORT"`
	Dbuser     string `required:"true" envconfig:"DB_USER"`
	Dbpassword string `required:"true" envconfig:"DB_PASS"`
	Dbname     string `required:"true" envconfig:"DB_NAME"`
	JwtKey     string `required:"true" envconfig:"JWT_KEY"`
	HashKey    string `required:"true" envconfig:"HASH_KEY"`
}

//Env holds application config variables
var Env envVars

// Tpl template
var Tpl *template.Template

func init() {

	Tpl = template.Must(template.ParseGlob("/"))

	//load .env file
	err := godotenv.Load("././.env")

	if err != nil {
		log.Println("Error loading .env file, falling back to cli passed env")
	}

	err = envconfig.Process("", &Env)

	if err != nil {
		log.Fatalln("Error loading environment variables", err)
	}

}
