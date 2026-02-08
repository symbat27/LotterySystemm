package storage

import (
	"LotterySystem/internal/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const ticketFile = "data/tickets.json"

type TicketRepository struct {
	mu sync.RWMutex
	db map[string]models.Ticket
}

func NewTicketRepository() *TicketRepository {
	r := &TicketRepository{
		db: make(map[string]models.Ticket),
	}
	r.load()
	return r
}

func (r *TicketRepository) load() {
	data, err := os.ReadFile(ticketFile)
	if err == nil {
		_ = json.Unmarshal(data, &r.db)
	}
}

func (r *TicketRepository) save() {
	data, _ := json.MarshalIndent(r.db, "", "  ")
	_ = os.WriteFile(ticketFile, data, 0644)
}

func (r *TicketRepository) Save(t models.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[t.ID] = t
	r.save()
	return nil
}

func (r *TicketRepository) Update(t models.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.db[t.ID]; !exists {
		return errors.New("ticket not found")
	}
	r.db[t.ID] = t
	r.save()
	return nil
}

func (r *TicketRepository) GetByID(id string) (models.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, exists := r.db[id]
	if !exists {
		return models.Ticket{}, errors.New("ticket not found")
	}
	return ticket, nil
}

func (r *TicketRepository) List() []models.Ticket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tickets := make([]models.Ticket, 0, len(r.db))
	for _, t := range r.db {
		tickets = append(tickets, t)
	}
	return tickets
}

func (r *TicketRepository) GetByUserID(userID string) []models.Ticket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := []models.Ticket{}
	for _, t := range r.db {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result
}

func (r *TicketRepository) GetByDrawID(drawID string) []models.Ticket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := []models.Ticket{}
	for _, t := range r.db {
		if t.DrawID == drawID {
			result = append(result, t)
		}
	}
	return result
}
