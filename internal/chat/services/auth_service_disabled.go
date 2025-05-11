package services

import "cu.ru/internal/chat/tokens"

type AuthServiceDisabled struct {
}

func NewAuthServiceDisabled() *AuthServiceDisabled {
	return &AuthServiceDisabled{}
}

func (a *AuthServiceDisabled) GenerateToken(login, password string) (string, error) {
	return tokens.GenerateToken(login, "doesn't matter")
}
