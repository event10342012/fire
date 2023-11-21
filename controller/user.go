package controller

import (
	"fire/model"
	"fire/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetUsers(c *gin.Context) {
	db := server.GetDB()

	if c.Query("skip") != "" {
		skip := cast.ToInt(c.Query("skip"))
		db = db.Offset(skip)
	}
	if c.Query("limit") != "" {
		limit := cast.ToInt(c.Query("limit"))
		db = db.Limit(limit)
	}

	var users []model.User
	db.Find(&users)

	for i := range users {
		// Clear the password for each user in the original slice
		users[i].Password = ""
	}

	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	db := server.GetDB()
	db.First(&user, id)
	c.JSON(200, user)
}

func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid payload",
		})
		return
	}

	user.Password = server.HashPassword(user.Password)
	db := server.GetDB()
	db.Create(&user)
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	db := server.GetDB()
	db.First(&user, id)

	var input model.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid payload",
		})
		return
	}

	if !user.CheckPassword(input.Password) {
		input.Password = server.HashPassword(input.Password)
	}

	db.Model(&user).Updates(&input)
	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	db := server.GetDB()
	db.First(&user, id)
	db.Delete(&user)
	c.JSON(200, gin.H{
		"message": "User deleted",
	})
}
