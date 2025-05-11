package services

type AuthService interface {
	GenerateToken(login, password string) (string, error)
}
