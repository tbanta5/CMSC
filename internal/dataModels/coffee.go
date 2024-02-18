package dataModels

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Metadata tags help us decode json and format for database.
type Coffee struct {
	ID          int     `db:"coffee_id" json:"id,omitempty"`
	Name        string  `db:"coffee_name" json:"name"`
	Description string  `db:"coffee_description" json:"description"`
	Price       float64 `db:"coffee_price" json:"price"`
	Caffeine    string  `db:"coffee_caffeine" json:"caffeine,omitempty"`
	Calories    int     `db:"coffee_calories" json:"calories,omitempty"`
}

// Slice of Coffees
// type Coffees []Coffee

func (c *Coffee) AddCoffee(ctx context.Context, db *pgxpool.Pool) error {
	return nil
}

func CoffeeList(ctx context.Context, db *pgxpool.Pool) ([]Coffee, error) {
	// Empty coffee slice
	coffees := []Coffee{}

	// Ensure db is connected.
	if err := db.Ping(ctx); err != nil {
		return coffees, err
	}
	// Prepare statement
	const stmt = `select * from coffee;`
	rows, err := db.Query(ctx, stmt)
	if err != nil {
		return coffees, err
	}
	// Scan returned coffees
	for rows.Next() {
		coffee := Coffee{}
		err := rows.Scan(
			&coffee.ID,
			&coffee.Name,
			&coffee.Description,
			&coffee.Price,
			&coffee.Caffeine,
			&coffee.Calories)
		if err != nil {
			return coffees, err
		}
		// Write to coffee slice
		coffees = append(coffees, coffee)
	}

	return coffees, nil
}

func CoffeeDetails(ctx context.Context, db *pgxpool.Pool, id int) (Coffee, error) {
	// Instantiate new coffee structure
	coffee := Coffee{}

	// Ensure db is connected.
	if err := db.Ping(ctx); err != nil {
		return Coffee{}, err
	}

	const stmt = `select * from coffee where coffee_id=$1;`
	row := db.QueryRow(ctx, stmt, id)

	if err := row.Scan(&coffee.ID, &coffee.Name, &coffee.Description, &coffee.Price, &coffee.Caffeine, &coffee.Calories); err != nil {
		return Coffee{}, err
	}

	return coffee, nil
}
