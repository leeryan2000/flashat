package repo

import "github.com/leeryan2000/flashat/model"

type UserRepository interface {
    CreateUser(user *model.User) error
    GetAllUsers() ([]model.User, error)
}