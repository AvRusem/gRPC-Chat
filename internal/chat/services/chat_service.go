package services

import "cu.ru/api/pb"

type ChatService interface {
	RegisterClient(clientID string, stream ChatStream) error
	UnregisterClient(clientID string) error
	BroadcastMessage(clientID, message string) error
	BanClient(clientID string) error
	ModerateMessage(clientID, message string) (bool, error)
}

type ChatStream interface {
	Send(*pb.ChatMessage) error
}
