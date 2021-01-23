package validator

import (
	"github.com/go-playground/validator/v10"
)

var (
	// use a single instance of Validate, it caches struct info
	Validate *validator.Validate
)

func init() {
	Validate = validator.New()
}
