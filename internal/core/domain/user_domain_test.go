package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserValidateBusinessRules(t *testing.T) {
	validUser := &User{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Age:      25,
		Password: "1234",
		Role:     "user",
	}
	err := validUser.ValidateBusinessRules()
	assert.NoError(t, err)

	invalideUser := &User{
		Name:     "John123",
		Email:    "john@yahoo.com",
		Age:      15,
		Password: "12",
		Role:     "admin",
	}

	err = invalideUser.ValidateBusinessRules()
	assert.Error(t, err)
}
