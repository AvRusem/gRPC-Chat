package services

import (
	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/repositories"
	"cu.ru/internal/chat/tokens"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceEnabled struct {
	authRepository repositories.AuthRepository
}

func NewAuthServiceEnabled(authRepo repositories.AuthRepository) *AuthServiceEnabled {
	return &AuthServiceEnabled{
		authRepository: authRepo,
	}
}

func (a *AuthServiceEnabled) GenerateToken(login, password string) (string, error) {
	user, err := a.authRepository.GetUser(login)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", appErrors.ErrNotAuthorized
	}

	return tokens.GenerateToken(login, user.Role)
}
