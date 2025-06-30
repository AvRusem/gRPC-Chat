package repositories

import (
	"os"
	"strings"
	"sync"
)

type ProfanityRepositoryInMemory struct {
	profanities map[string]struct{}
	mu          sync.RWMutex
}

func (r *ProfanityRepositoryInMemory) loadProfanities(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	profanities := make([]string, 0, len(lines))
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			profanities = append(profanities, strings.ToLower(word))
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, word := range profanities {
		r.profanities[word] = struct{}{}
	}
}

func NewProfanityRepositoryInMemory(file string) *ProfanityRepositoryInMemory {
	profanities := make(map[string]struct{})
	repo := &ProfanityRepositoryInMemory{
		profanities: profanities,
	}
	repo.loadProfanities(file)
	return repo
}

func (r *ProfanityRepositoryInMemory) ContainsProfanity(text string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	text = strings.ToLower(strings.Join(strings.Fields(text), " "))
	for word := range r.profanities {
		if strings.Contains(text, word) {
			return true
		}
	}

	return false
}
