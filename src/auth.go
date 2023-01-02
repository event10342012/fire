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
	Password string `json:"password"`
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/fire")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (user User) Save() error {
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

	row := db.QueryRow("select username, password from user where id = ?", id)

	var username string
	var password string

	err := row.Scan(&username, &password)
	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}

	return user, err
}

func GetUserByUsername(username string) (User, error) {
	db := getDB()
	defer db.Close()

	row := db.QueryRow("select id, password from user where username = ?", username)

	var id int
	var password string

	err := row.Scan(&id, &password)

	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}

	return user, err
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	loginUser := User{
		Username: username,
		Password: password,
	}

	userDB, err := GetUserByUsername(username)

	if err != nil {
		c.String(http.StatusNotFound, "User: %s not found", username)
		return
	}

	if userDB.Password != loginUser.Password {
		c.String(http.StatusUnauthorized, "Password is not correct")
		return
	}

	c.JSON(http.StatusOK, "Login success")
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if password != confirmPassword {
		c.String(http.StatusNotFound, "password and confirm password dose not match")
		return
	}

	user := User{
		Username: username,
		Password: password,
	}
	err := user.Save()
	if err != nil {
		c.String(http.StatusNotFound, "User: %s already exists")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("User: %s register success", username))
}

func AddAuthRoutes(rg *gin.Engine) {
	auth := rg.Group("/auth")

	auth.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	auth.POST("/login", login)

	auth.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	auth.POST("/register", register)
}
