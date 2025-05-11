package repositories

import "cu.ru/internal/chat/models"

type AuthRepository interface {
	GetPasswordHash(login string) (models.AuthUser, error)
}
