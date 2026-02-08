package storage

import (
	"LotterySystem/internal/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const prizeFile = "data/prizes.json"

type PrizeRepository struct {
	mu sync.RWMutex
	db map[string]models.Prize
}

func NewPrizeRepository() *PrizeRepository {
	r := &PrizeRepository{
		db: make(map[string]models.Prize),
	}
	r.load()
	return r
}

func (r *PrizeRepository) load() {
	data, err := os.ReadFile(prizeFile)
	if err == nil {
		_ = json.Unmarshal(data, &r.db)
	}
}

func (r *PrizeRepository) save() {
	data, _ := json.MarshalIndent(r.db, "", "  ")
	_ = os.WriteFile(prizeFile, data, 0644)
}

func (r *PrizeRepository) Save(p models.Prize) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[p.ID] = p
	r.save()
	return nil
}

func (r *PrizeRepository) GetByID(id string) (models.Prize, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prize, exists := r.db[id]
	if !exists {
		return models.Prize{}, errors.New("prize not found")
	}
	return prize, nil
}

func (r *PrizeRepository) GetByTicketID(ticketID string) (models.Prize, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, prize := range r.db {
		if prize.TicketID == ticketID {
			return prize, nil
		}
	}
	return models.Prize{}, errors.New("prize not found for ticket")
}

func (r *PrizeRepository) List() []models.Prize {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prizes := make([]models.Prize, 0, len(r.db))
	for _, p := range r.db {
		prizes = append(prizes, p)
	}
	return prizes
}
