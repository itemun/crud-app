package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/itemun/crud-app/internal/domain"
)

type Cars struct {
	db *sql.DB
}

func NewCars(db *sql.DB) *Cars {
	return &Cars{db}
}

func (c *Cars) Create(ctx context.Context, car domain.Car) error {
	_, err := c.db.Exec("INSERT INTO cars (model, name, production_date, hp) values ($1, $2, $3, $4)",
		car.Model, car.Name, car.ProductionDate, car.HP)

	return err
}

func (c *Cars) GetByID(ctx context.Context, id int64) (domain.Car, error) {
	var car domain.Car
	err := c.db.QueryRow("SELECT id, model, name, production_date, hp FROM cars WHERE id=$1", id).
		Scan(&car.ID, &car.Model, &car.Name, &car.ProductionDate, &car.HP)
	if err == sql.ErrNoRows {
		return car, domain.ErrCarNotFound
	}

	return car, err
}

func (c *Cars) GetAll(ctx context.Context) ([]domain.Car, error) {
	rows, err := c.db.Query("SELECT id, model, name, production_date, hp FROM cars")
	if err != nil {
		return nil, err
	}

	cars := make([]domain.Car, 0)
	for rows.Next() {
		var car domain.Car
		if err := rows.Scan(&car.ID, &car.Model, &car.Name, &car.ProductionDate, &car.HP); err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}

	return cars, rows.Err()
}

func (c *Cars) Delete(ctx context.Context, id int64) error {
	_, err := c.db.Exec("DELETE FROM cars WHERE id=$1", id)

	return err
}

func (c *Cars) Update(ctx context.Context, id int64, inp domain.UpdateCarInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Model != nil {
		setValues = append(setValues, fmt.Sprintf("model=$%d", argId))
		args = append(args, *inp.Model)
		argId++
	}

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *inp.Name)
		argId++
	}

	if inp.ProductionDate != nil {
		setValues = append(setValues, fmt.Sprintf("production_date=$%d", argId))
		args = append(args, *inp.ProductionDate)
		argId++
	}

	if inp.HP != nil {
		setValues = append(setValues, fmt.Sprintf("hp=$%d", argId))
		args = append(args, *inp.HP)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE cars SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := c.db.Exec(query, args...)
	return err
}
