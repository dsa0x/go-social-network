package model

import (
	"log"

	"github.com/dsa0x/go-social-network/common"
	"github.com/dsa0x/go-social-network/internal/config"
	"github.com/jinzhu/gorm"
)

// Db database
var Db = config.DB()

func init() {
	// defer db.Close()
	// Db.CreateTable(&User{})
	Db.DropTableIfExists(&User{})
	Db.CreateTable(&User{})
	Db.AutoMigrate(&User{})
}

// User struct declaration
type User struct {
	gorm.Model
	UserName        string `json:"username"  binding:"required" gorm:"unique;not null"`
	Email           string `json:"email"  binding:"required" gorm:"unique;not null"`
	Friends         []User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
	Password        string `json:"password,omitempty"  binding:"required"`
	ConfirmPassword string `json:"confirmPassword,omitempty"  binding:"required" sql:"-"`
}

func FetchUsers(users *[]User) error {
	allUsers := Db.Model(&User{}).Order("user_name asc").Find(&users)
	if allUsers.Error != nil {
		log.Println(allUsers.Error)
		return allUsers.Error
	}

	return nil
}

// FindOne finds one user by emails
func FindOne(email string) (*User, error) {
	user := &User{}

	if err := Db.Where("email = ?", email).First(user).Error; err != nil {
		// var resp := map[string]interface{} {"status": false, "message": "invalid "}
		return user, err
	}

	return user, nil

}

// FindByID finds one user by id
func FindByID(ID uint) (*User, error) {
	user := &User{}

	if err := Db.First(user, ID).Error; err != nil {
		return user, err
	}

	return user, nil

}

// CreateUser creates a new user
func CreateUser(user User) (uint, error) {
	pass, err := common.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		return user.ID, err
	}
	user.Password = string(pass)
	createdUser := Db.Create(&user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		log.Println(errMessage)
		return user.ID, errMessage
	}
	return user.ID, nil
}

func Follow(userID uint, friendID uint) (uint, error) {

	user := User{}
	friend, err := FindByID(friendID)
	if err != nil {
		log.Println(err)
		return userID, err
	}
	Db.Preload("Friends").First(&user, "id = ?", userID)
	Db.Model(&user).Association("Friends").Append(friend)
	return userID, nil
}

func Unfollow(userID uint, friendID uint) (uint, error) {

	user := User{}
	friend, err := FindByID(friendID)
	if err != nil {
		log.Println(err)
		return userID, err
	}
	Db.Preload("Friends").First(&user, "id = ?", userID)
	Db.Model(&user).Association("Friends").Delete(friend)
	return userID, nil
}

func IsFollower(userID uint, friendID uint) (bool, error) {

	friend := User{}
	user := User{}
	friend.ID = friendID
	user.ID = userID
	Db.Model(&user).Association("Friends").Find(&friend)
	if friend.Email == "" {
		return false, nil
	}
	return true, nil
}
func CountFollowings(userID uint) (int, error) {
	user := User{}
	user.ID = userID
	count := Db.Model(&user).Association("Friends").Count()
	return count, nil
}
func CountFollowers(userID uint) (int, error) {
	count := 0
	Db.Table("friendships").Select("friend_id").Where("friend_id = ?", userID).Count(&count)
	return count, nil
}
