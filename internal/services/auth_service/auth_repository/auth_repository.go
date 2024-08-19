package authrepository

import (
	database "github.com/raganrrlaw/server/internal/db"
)

type AuthRepository interface {
  GetToken(string) (string, error) // get the token
  GetAssociatedToken(string) (string, error) // get the token
  AddToken(string, string) error // crate a token -> associate id (user id), token
  RemoveToken(string) error // remove token -> using the token
}

type AuthRepo struct {
	storage *database.Storage
}

func NewAuthRepo(storage *database.Storage) *AuthRepo {
  return &AuthRepo{
    storage: storage,
  }
}

func (ar *AuthRepo) GetToken(id string) (string, error) {
  return "", nil
}

func (ar *AuthRepo) AddToken(assocId string, token string) error {
  return nil
}

func (ar *AuthRepo) GetAssociatedToken(assocId string) (string, error) {
  return "", nil
}

func (ar *AuthRepo) RemoveToken(token string) error {
  return nil
}

