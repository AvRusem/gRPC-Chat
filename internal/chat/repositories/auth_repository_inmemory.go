package repositories

import (
	"bufio"
	"os"
	"strings"
	"sync"

	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/models"
)

type AuthRepositoryInMemory struct {
	users map[string]models.AuthUser
	mu    sync.RWMutex
}

func (r *AuthRepositoryInMemory) loadUsers(filePath string) {
	// file's row is login:hashed_password:role
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			continue
		}
		login := parts[0]
		passwordHash := parts[1]
		role := parts[2]

		r.mu.Lock()
		r.users[login] = models.AuthUser{
			Login:        login,
			PasswordHash: passwordHash,
			Role:         role,
		}
		r.mu.Unlock()
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func NewAuthRepositoryInMemory(authStorage string) *AuthRepositoryInMemory {
	users := make(map[string]models.AuthUser)
	repo := &AuthRepositoryInMemory{
		users: users,
	}
	repo.loadUsers(authStorage)
	return repo
}

func (r *AuthRepositoryInMemory) GetUser(login string) (models.AuthUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[login]
	if !exists {
		return models.AuthUser{}, appErrors.ErrNotFound
	}
	return user, nil
}
