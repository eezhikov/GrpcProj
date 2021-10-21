package models

type UserList struct {
	User *User
}
type User struct {
	Id    int
	Password string
	Login string
	Name  string
	Phone string
}