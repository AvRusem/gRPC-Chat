package repositories

import (
	"sync"

	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/models"
)

type ChatRepositoryInMemory struct {
	users map[string]models.ChatUser
	mu    sync.RWMutex
}

func NewChatRepositoryInMemory() *ChatRepositoryInMemory {
	users := make(map[string]models.ChatUser)
	repo := &ChatRepositoryInMemory{
		users: users,
	}
	return repo
}

func (r *ChatRepositoryInMemory) AddUser(login string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[login]; exists {
		return nil
	}
	r.users[login] = models.ChatUser{
		Login:   login,
		Penalty: 0,
		Banned:  false,
	}

	return nil
}

func (r *ChatRepositoryInMemory) PunishUser(login string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, exists := r.users[login]; exists {
		user.Penalty++
		if user.Penalty >= 3 {
			user.Banned = true
		}
		r.users[login] = user
		return nil
	}
	return appErrors.ErrNotFound
}

func (r *ChatRepositoryInMemory) BanUser(login string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, exists := r.users[login]; exists {
		user.Banned = true
		r.users[login] = user
		return nil
	}
	return appErrors.ErrNotFound
}

func (r *ChatRepositoryInMemory) IsBanned(login string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, exists := r.users[login]; exists {
		return user.Banned, nil
	}
	return false, appErrors.ErrNotFound
}
