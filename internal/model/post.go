package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model `json:"-"`
	Title      string
	Content    string
	CreatedBy  uint
	Author     string
}

func init() {
	// Db.CreateTable(&Post{})
}

// CreatePost create post
func CreatePost(post Post) (uint, error) {
	createdPost := Db.Create(&post)
	if createdPost.Error != nil {
		log.Println(createdPost.Error)
		return post.ID, createdPost.Error
	}
	return post.ID, nil
}

// DeletePost delete post
func DeletePost(postID uint, userID uint) (uint, error) {
	deletedPost := Db.Unscoped().Delete(Post{}, "ID = ? and created_by = ?", postID, userID)
	if deletedPost.Error != nil {
		log.Println(deletedPost.Error)
		return postID, deletedPost.Error
	}

	return postID, nil
}

// UpdatePost updates the post
func UpdatePost(post Post, postID uint, userID uint) (Post, error) {
	oldPost := Post{}
	updatedPost := Db.First(&oldPost, postID)
	oldPost.Title = post.Title
	oldPost.Content = post.Content
	updatedPost = Db.Save(&oldPost)
	if updatedPost.Error != nil {
		log.Println(updatedPost.Error)
		return oldPost, updatedPost.Error
	}

	return oldPost, nil
}

// FetchPosts fetch post
func FetchPosts(posts *[]Post) error {
	allPosts := Db.Model(&Post{}).Order("created_at desc").Find(&posts)
	if allPosts.Error != nil {
		log.Println(allPosts.Error)
		return allPosts.Error
	}

	return nil
}

func FetchUserPosts(posts *[]Post, userID uint) error {
	allPosts := Db.Model(&Post{}).Order("created_at desc").Find(&posts, "created_by = ?", userID)
	if allPosts.Error != nil {
		log.Println(allPosts.Error)
		return allPosts.Error
	}

	return nil
}
func CountUserPosts(userID uint) (int, error) {
	var count int
	allPosts := Db.Model(&Post{}).Where("created_by = ?", userID).Count(&count)
	if allPosts.Error != nil {
		log.Println(allPosts.Error)
		return count, allPosts.Error
	}

	return count, nil
}
