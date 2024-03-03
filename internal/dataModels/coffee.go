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

func AddCoffee(ctx context.Context, db *pgxpool.Pool, c Coffee) (int, error) {
	const stmt = `
	INSERT INTO coffee 
	(coffee_name, coffee_description, coffee_price, coffee_caffeine, coffee_calories) 
	VALUES ($1, $2, $3, $4, $5) RETURNING coffee_id`
	err := db.QueryRow(ctx, stmt, c.Name, c.Description, c.Price, c.Caffeine, c.Calories).Scan(&c.ID)
	if err != nil {
		return 0, err
	}
	return c.ID, nil
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

func DeleteCoffee(ctx context.Context, db *pgxpool.Pool, id int) error {

	const stmt = `DELETE from coffee where coffee_id=$1;`
	_, err := db.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCoffee(ctx context.Context, db *pgxpool.Pool, id int, coffee Coffee) error {
	const stmt = `UPDATE coffee set
	coffee_name = $1, 
	coffee_description = $2, 
	coffee_price = $3,
	coffee_caffeine = $4, 
	coffee_calories = $5 
	WHERE coffee_id = $6
	`
	args := []any{coffee.Name, coffee.Description, coffee.Price, coffee.Caffeine, coffee.Calories, id}

	_, err := db.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}
	return nil
}

func CheckCoffeeExists(ctx context.Context, db *pgxpool.Pool, id int) bool {
	const stmt = `select coffee_id from coffee where coffee_id=$1`
	_, err := db.Exec(ctx, stmt, id)
	// This will be false or true depending.
	return err != nil
}
