package services

import (
	"log"
	"sync"

	"cu.ru/api/pb"
	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/repositories"
)

type ChatServiceDefault struct {
	clients             map[string]ChatStream
	mu                  sync.RWMutex
	profanityRepository repositories.ProfanityRepository
	chatRepository      repositories.ChatRepository
}

func NewChatServiceDefault(
	profanityRepository repositories.ProfanityRepository,
	chatRepository repositories.ChatRepository,
) *ChatServiceDefault {
	return &ChatServiceDefault{
		clients:             make(map[string]ChatStream),
		profanityRepository: profanityRepository,
		chatRepository:      chatRepository,
	}
}

func (s *ChatServiceDefault) RegisterClient(clientID string, stream ChatStream) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[clientID]; exists {
		return appErrors.AlreadyExistsError
	}
	s.clients[clientID] = stream
	return s.chatRepository.AddUser(clientID)
}

func (s *ChatServiceDefault) UnregisterClient(clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[clientID]; !exists {
		return appErrors.NotFoundError
	}
	delete(s.clients, clientID)
	return nil
}

func (s *ChatServiceDefault) BroadcastMessage(clientID, message string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.clients[clientID]; !exists {
		return appErrors.NotFoundError
	}

	banned, err := s.chatRepository.IsBanned(clientID)
	if err != nil {
		return err
	}
	if banned {
		return appErrors.BannedError
	}

	for id, stream := range s.clients {
		err := stream.Send(&pb.ChatMessage{Login: clientID, Text: message})
		if err != nil {
			// Need to ignore the error
			// as the client may have disconnected
			// and we don't want to crash the server
			// or stop broadcasting to other clients
			log.Printf("failed to send message to client %s: %v", id, err)
			continue
		}
	}
	return nil
}

func (s *ChatServiceDefault) BanClient(clientID string) error {
	err := s.chatRepository.BanUser(clientID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChatServiceDefault) ModerateMessage(clientID, message string) (bool, error) {
	banned, err := s.chatRepository.IsBanned(clientID)
	if err != nil {
		return true, err
	}
	if banned {
		return true, appErrors.BannedError
	}

	if s.profanityRepository.ContainsProfanity(message) {
		err := s.chatRepository.PunishUser(clientID)
		if err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}
