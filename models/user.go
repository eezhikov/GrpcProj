package models

type UserList struct {
	User []*UserInfo
}
type UserInfo struct {
	Id       int
	Password string
	Login    string
	Name     string
	Phone    string
}
