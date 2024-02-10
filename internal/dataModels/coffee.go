package dataModels

import "context"

// Metadata tags help us decode json and format for database.
type Coffee struct {
	ID          int    `db:"id" json:"id,omitempty"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

// Slice of Coffees
type Coffees []Coffee

func (c *Coffee) AddCoffee(ctx context.Context) {

}

func CoffeeList(ctx context.Context) {

}
