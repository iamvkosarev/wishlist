package storage

import (
	"github.com/iamvkosarev/wishlist/back/internal/model"
)

type LocalStorage struct {
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (l LocalStorage) GetUser(email string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}
