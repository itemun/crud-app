package service

import (
	"context"
	"time"

	"github.com/itemun/crud-app/internal/domain"
)

type CarsRepository interface {
	Create(ctx context.Context, car domain.Car) error
	GetByID(ctx context.Context, id int64) (domain.Car, error)
	GetAll(ctx context.Context) ([]domain.Car, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateCarInput) error
}

type Cars struct {
	repo CarsRepository
}

func NewCars(repo CarsRepository) *Cars {
	return &Cars{
		repo: repo,
	}
}

func (c *Cars) Create(ctx context.Context, car domain.Car) error {
	if car.ProductionDate.IsZero() {
		car.ProductionDate = time.Now()
	}

	return c.repo.Create(ctx, car)
}

func (c *Cars) GetByID(ctx context.Context, id int64) (domain.Car, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *Cars) GetAll(ctx context.Context) ([]domain.Car, error) {
	return c.repo.GetAll(ctx)
}

func (c *Cars) Delete(ctx context.Context, id int64) error {
	return c.repo.Delete(ctx, id)
}

func (c *Cars) Update(ctx context.Context, id int64, inp domain.UpdateCarInput) error {
	return c.repo.Update(ctx, id, inp)
}
