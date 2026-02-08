package models

import "time"

type Ticket struct {
	ID        string    json:"id"
	UserID    string    json:"user_id"
	DrawID    string    json:"draw_id"
	Numbers   []int     json:"numbers"
	Matches   int       json:"matches"
	PrizeID   string    json:"prize_id,omitempty"
	CreatedAt time.Time json:"created_at"
}