package models

type AuthUser struct {
	Login        string
	Role         string
	PasswordHash string
}

type ChatUser struct {
	Login   string
	Penalty int
	Banned  bool
}
