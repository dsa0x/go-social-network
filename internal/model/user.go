package model

import (
	"github.com/dsa0x/social-network/model"
	"github.com/jinzhu/gorm"
)

// User struct declaration
type User struct {
	gorm.Model
	Name     string `json:"name"  binding:"required"`
	Email    string `json:"email"  binding:"required"`
	Gender   string `json:"gender" `
	Password string `json:"password"  binding:"required"`
}

func FindOne(email string, password string) {
	user := &model.User
}
