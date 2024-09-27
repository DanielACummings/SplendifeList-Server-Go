package models

import (
	"time"
)

type Item struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CrossedOut bool      `json:"crossed_out"`
	User       uint      `json:"user"`
}
