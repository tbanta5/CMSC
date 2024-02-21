package validation

import (
	"cmsc.group2.coffee-api/internal/dataModels"
	"github.com/go-playground/validator/v10"
)

// Initialize a validator instance
var validate = validator.New()

func ValidateCoffee(coffee *dataModels.Coffee) error {
	// Perform validation and return any errors
	return validate.Struct(coffee)
}
