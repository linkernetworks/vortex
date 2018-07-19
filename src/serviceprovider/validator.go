package serviceprovider

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func checkNameValidation(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])`)
	return re.MatchString(fl.Field().String())
}
