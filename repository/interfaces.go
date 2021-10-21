package repository

import "UserGrpcProj/models"

type UserRepository interface {
	AddUser(user *models.User) (bool, error)
	RemoveUser(int) (bool, error)
	UserList(string) (*models.UserList, error)
	IsSetUser(string, int) bool
}
