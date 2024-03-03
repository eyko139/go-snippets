package models

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
    ID string
    Name string
    Email string
    HashedPassword []byte
    Create time.Time
}

type UserModel struct {
    DbClient *mongo.Client
}

type UserModelInterface interface {
    Insert(name, email, password string) error
    Authenticate(email, password string) (int, error)
    Exists(id int) (bool, error)
}

func (um *UserModel) Insert(name, email, password string) error {
    return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
    return 0, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
    return false, nil 
}
