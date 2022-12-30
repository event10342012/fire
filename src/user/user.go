package user

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:seafin@tcp(127.0.0.1)/fire")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (user User) Save() {
	db := getDB()
	defer db.Close()

	insertStat, err := db.Prepare("INSERT INTO user(username, email, hash_password) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertStat.Close()
	_, err = insertStat.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		panic(err.Error())
	}
}

func GetUser(id int) {
	db := getDB()
	defer db.Close()

	result, err := db.Query("select * from user where id = ?", id)
	if err != nil {
		panic(err.Error())
	}

	columns, err := result.Columns()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(columns)
}
