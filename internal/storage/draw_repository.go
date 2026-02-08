package storage

import (
	"LotterySystem/internal/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const drawFile = "data/draws.json"

type DrawRepository struct {
	mu sync.RWMutex
	db map[string]models.Draw
}

func NewDrawRepository() *DrawRepository {
	r := &DrawRepository{
		db: make(map[string]models.Draw),
	}
	r.load()
	return r
}

func (r *DrawRepository) load() {
	data, err := os.ReadFile(drawFile)
	if err == nil {
		_ = json.Unmarshal(data, &r.db)
	}
}

func (r *DrawRepository) save() {
	data, _ := json.MarshalIndent(r.db, "", "  ")
	_ = os.WriteFile(drawFile, data, 0644)
}

func (r *DrawRepository) Save(d models.Draw) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[d.ID] = d
	r.save()
	return nil
}

func (r *DrawRepository) Update(d models.Draw) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[d.ID] = d
	r.save()
	return nil
}

func (r *DrawRepository) GetByID(id string) (models.Draw, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.db[id]
	if !ok {
		return models.Draw{}, errors.New("draw not found")
	}
	return d, nil
}

func (r *DrawRepository) GetPending() (models.Draw, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, d := range r.db {
		if d.Status == "pending" {
			return d, nil
		}
	}
	return models.Draw{}, errors.New("no pending draw")
}

func (r *DrawRepository) List() []models.Draw {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []models.Draw{}
	for _, d := range r.db {
		res = append(res, d)
	}
	return res
}
