package serviceprovider

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("k8sname", checkNameValidation)
}

func TestCheckNameValidation(t *testing.T) {
	name := "awesome"
	err := validate.Var(name, "required,k8sname")
	assert.NoError(t, err)

	name = "123abc"
	err = validate.Var(name, "required,k8sname")
	assert.NoError(t, err)
}

func TestCheckNameValidationFail(t *testing.T) {
	name := "~!@#$%^&*()"
	err := validate.Var(name, "required,k8sname")
	assert.Error(t, err)
}
