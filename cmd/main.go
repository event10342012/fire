package main

import (
	"fire/src/user"
)

func main() {
	//u := user.User{
	//	Username: "leochen",
	//	Email:    "leochen@gmail.com",
	//	Password: "12345678",
	//}
	//
	//u.Save()

	user.GetUser(3)
}
