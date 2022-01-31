package controllers

import (
	"fmt"
	"time"

	"com.uf/src/models"
	"com.uf/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func DeletePost(c *gin.Context) {
	id := c.Params.ByName("id")
	var post models.UserPost
	d := utils.DB.Where("post_id = ?", id).Delete(&post)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func PostComment(c *gin.Context) {
	var comment models.Comment
	c.BindJSON(&comment)
	var post models.UserPost
	result := utils.DB.Where("post_id = ?", comment.PostID).First(&post)

	if result != nil && result.RowsAffected == 1 {
		comment.CreatedAt = time.Now().Unix()
		utils.DB.Create(&comment)
		c.JSON(200, comment)
	} else {
		c.JSON(400, gin.H{"error": "Unable to add comment"})
	}
}

func UpdateLikes(c *gin.Context) {
	var post models.UserPost
	post_id := c.Params.ByName("post_id")

	liked := c.Params.ByName("liked")

	result := utils.DB.Where("post_id = ?", post_id).First(&post)

	if result != nil && result.RowsAffected == 1 {
		if liked == "true" {
			post.Likes = post.Likes + 1
		} else {
			post.Likes = post.Likes - 1
		}
		utils.DB.Save(&post)
		c.JSON(200, gin.H{"message": "likes updated"})
	} else {
		c.JSON(400, gin.H{"error": "Unable to update likes"})
	}
}

func UpdatePost(c *gin.Context) {
	var post models.UserPost
	id := c.Params.ByName("id")
	if err := utils.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&post)
	post.UpdatedAt = time.Now().Unix()
	utils.DB.Save(&post)
	c.JSON(200, post)
}

func CreatePost(c *gin.Context) {
	var post models.UserPost
	c.BindJSON(&post)
	post.CreatedAt = time.Now().Unix()
	utils.DB.Create(&post)
	c.JSON(200, post)
}

func GetPost(c *gin.Context) {
	id := c.Params.ByName("id")
	var post models.UserPost
	if err := utils.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, post)
	}
}

func GetPosts(c *gin.Context) {
	var userposts []models.UserPost
	result := utils.DB.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at desc")
		return db
	}).Order("created_at desc").Find(&userposts)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Unable to retrieve feed posts"})
	} else {
		c.JSON(200, userposts)
	}
}
