package models

import "time"

type Draw struct {
	ID             string    json:"id"
	WinningNumbers []int     json:"winning_numbers"
	Status         string    json:"status" // "pending", "completed"
	DrawDate       time.Time json:"draw_date"
	CreatedAt      time.Time json:"created_at"
}