package controller

import (
	"fire/model"
	"fire/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type userBase struct {
	Email       string    `json:"email"`
	GivenName   string    `json:"given_name"`
	FamilyName  string    `json:"family_name"`
	Birthdate   time.Time `json:"birthdate"`
	IsSuperUser bool      `json:"is_super_user"`
	IsActive    bool      `json:"is_active"`
}

type userCreate struct {
	userBase
	Password string
}

type userUpdate struct {
	userCreate
}

func GetUserByID(c *gin.Context) {
	db := server.GetDB()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	user := model.User{}
	err = user.Read(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, "Can not find this user with id: "+idParam)
	}

	userSchema := userBase{
		Email:       user.Email,
		GivenName:   user.GivenName,
		FamilyName:  user.FamilyName,
		IsSuperUser: user.IsSuperUser,
		IsActive:    user.IsActive,
	}

	c.JSON(http.StatusOK, userSchema)
}

func PostUser(c *gin.Context) {
	var input userCreate
	err := c.BindJSON(&input)
	if err != nil {
		return
	}
	hashPassword := server.HashPassword(input.Password)

	user := model.User{
		Email:       input.Email,
		Password:    hashPassword,
		GivenName:   input.GivenName,
		FamilyName:  input.FamilyName,
		Picture:     "",
		Locale:      "",
		GoogleId:    "",
		IsSuperUser: false,
		IsActive:    true,
	}

	if err != nil {
		return
	}
	err = user.Create(server.GetDB())
	if err != nil {
		return
	}
}

func PutUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	var input userUpdate
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, "data is invalid")
	}

	db := server.GetDB()

	user := model.User{}
	err = user.Read(db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "email not found")
	}

	if input.Password != "" {
		if !user.CheckPassword(input.Password) {
			user.Password = server.HashPassword(input.Password)
		}
	}

	user.GivenName = input.GivenName
	user.FamilyName = input.FamilyName
	user.Birthdate = input.Birthdate
	err = user.Update(db)
	if err != nil {
		c.JSON(http.StatusBadRequest, "data is invalid")
	}
	c.JSON(http.StatusOK, "update user success")
}
