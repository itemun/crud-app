package domain

import (
	"errors"
	"time"
)

var (
	ErrCarNotFound = errors.New("car not found")
)

type Car struct {
	ID             int64     `json:"id"`
	Model          string    `json:"model"`
	Name           string    `json:"name"`
	ProductionDate time.Time `json:"production_date"`
	HP             int       `json:"hp"`
}

type UpdateCarInput struct {
	Model          *string    `json:"model"`
	Name           *string    `json:"name"`
	ProductionDate *time.Time `json:"production_date"`
	HP             *int       `json:"hp"`
}
