package src

import (
	"dbapp/utils"

	"github.com/go-playground/validator/v10"
)

var validateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// check supported currency
		return utils.IsSupportedCurrency(currency)
	}

	return false
}
