package repositories

import "cu.ru/internal/chat/models"

type AuthRepository interface {
	GetUser(login string) (models.AuthUser, error)
}
