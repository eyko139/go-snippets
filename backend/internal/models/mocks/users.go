package mocks

import (
	"github.com/eyko139/go-snippets/internal/models"
	"time"
    "errors"
)

var MockUser = &models.User{
    ID: "1",
    Name: "MockUser",
    Email: "mock@mock.de",
    HashedPassword: []byte("123"),
	Create: time.Now(),
}

type UserModel struct {}

// var MockUserModel = &models.UserModel{}

func (mum *UserModel) Insert(name, email, password string) error {
    return  nil
}

func (mum *UserModel) Authenticate(email, password string) (int, error) {
    return 1, nil
}
func (mum *UserModel) Exists(id int) (bool, error) {
    if id == 1 {
        return true, nil
    }
    return false, errors.New("User doesnt exist")
}
