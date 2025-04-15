package models

import (
	"time"
)

type Lab struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Complete  bool      `json:"complete"`
}
