package repositories

type ProfanityRepository interface {
	ContainsProfanity(text string) bool
}
