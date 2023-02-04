package service

import (
	"errors"

	"Forum"
	"Forum/pkg/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user Forum.User) error {
	if err := checkEmail(user.Email); err != nil {
		return err
	}
	if err := checkUsername(user.Username); err != nil {
		return err
	}
	ok := checkPassword(user.Password)
	if !ok {
		return errors.New("password is not valid")
	}
	_, err := a.repo.CheckInvalid(user.Email, user.Username)
	if err == nil {
		return ErrCheckInvalid
	}
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return err
	}
	return a.repo.CreateUser(user)
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return string(bytes), err
	}
	return string(bytes), nil
}
