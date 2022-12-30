package user

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id       int    `json:"id"`
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

func (user User) Save() {
	db := getDB()
	defer db.Close()

	insertStat, err := db.Prepare("INSERT INTO user(username, password) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertStat.Close()
	_, err = insertStat.Exec(user.Username, user.Password)
	if err != nil {
		panic(err.Error())
	}
}

func GetUser(id int) {
	db := getDB()
	defer db.Close()

	rows, err := db.Query("select * from user where id = ?", id)
	if err != nil {
		panic(err.Error())
	}

	results := make([]User, 10)

	for rows.Next() {
		var id int
		var username string
		var password string

		rows.Scan(&id, &username, &password)
		results = append(results, User{
			id:       id,
			Username: username,
			Password: password,
		})

		fmt.Println(results)
	}
}
