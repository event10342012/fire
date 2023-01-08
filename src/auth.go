package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type UserDB struct {
	User
	Password string `json:"password"`
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/fire")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (user UserDB) Save() error {
	db := getDB()
	defer db.Close()

	insertStat, err := db.Prepare("INSERT INTO user(username, password) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertStat.Close()
	_, err = insertStat.Exec(user.Username, user.Password)
	return err
}

func GetUserByID(id int) (User, error) {
	db := getDB()
	defer db.Close()

	row := db.QueryRow("select username from user where id = ?", id)

	var username string

	err := row.Scan(&username)
	user := User{
		Id:       id,
		Username: username,
	}

	return user, err
}

func getUserByUsername(username string) (UserDB, error) {
	db := getDB()
	defer db.Close()

	row := db.QueryRow("select id, password from user where username = ?", username)

	var id int
	var password string

	err := row.Scan(&id, &password)

	user := UserDB{
		User: User{
			Id:       id,
			Username: username,
		},
		Password: password,
	}

	return user, err
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	loginUser := UserDB{
		User: User{
			Username: username,
		},
		Password: password,
	}

	userDB, err := getUserByUsername(username)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("User: %s not found", username)})
		return
	}

	if userDB.Password != loginUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password is not correct"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login success"})
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if password != confirmPassword {
		c.JSON(http.StatusNotFound, gin.H{"message": "password and confirm password dose not match"})
		return
	}

	user := UserDB{
		User:     User{Username: username},
		Password: password,
	}
	err := user.Save()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("User: %s already exists", username)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User: %s register success", username)})
}

func AddAuthRoutes(rg *gin.Engine) {
	auth := rg.Group("/auth")

	auth.POST("/login", login)
	auth.POST("/register", register)
}
