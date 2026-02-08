package storage

import (
	"LotterySystem/internal/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const userFile = "data/users.json"

type UserRepository struct {
	mu sync.RWMutex
	db map[string]models.User
}

func NewUserRepository() *UserRepository {
	r := &UserRepository{
		db: make(map[string]models.User),
	}
	r.load()
	return r
}

func (r *UserRepository) load() {
	data, err := os.ReadFile(userFile)
	if err == nil {
		_ = json.Unmarshal(data, &r.db)
	}
}

func (r *UserRepository) save() {
	data, _ := json.MarshalIndent(r.db, "", "  ")
	_ = os.WriteFile(userFile, data, 0644)
}

func (r *UserRepository) Save(u models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[u.ID] = u
	r.save()
	return nil
}

func (r *UserRepository) Update(u models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.db[u.ID]; !ok {
		return errors.New("user not found")
	}
	r.db[u.ID] = u
	r.save()
	return nil
}

func (r *UserRepository) GetByID(id string) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.db[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(username string) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.db {
		if u.Username == username {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (r *UserRepository) List() []models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []models.User{}
	for _, u := range r.db {
		res = append(res, u)
	}
	return res
}
