package repositories

type ChatRepository interface {
	AddUser(login string) error
	PunishUser(login string) error
	BanUser(login string) error
	IsBanned(login string) (bool, error)
}
