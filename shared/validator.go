package shared

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var v *validator.Validate

// GetValidator is responsible for returning a single instance of the validator.
func GetValidator() *validator.Validate {
	once.Do(func() {
		log.Info().Msg("Validator initialized.")
		v = validator.New()
	})

	return v
}

// AlphaSpace is function validation string contains only string & space
func AlphaSpace(fl validator.FieldLevel) bool {
	matched, err := regexp.Match("^[a-zA-Z\\s]+$", []byte(fl.Field().String()))
	if matched {
		return true
	}
	if err != nil {
		return false
	}

	return false
}

// AlphaNumSpace is function validation string contains only string, numeric & space
func AlphaNumSpace(fl validator.FieldLevel) bool {
	matched, err := regexp.Match("^[a-zA-Z0-9\\s]+$", []byte(fl.Field().String()))
	if matched {
		return true
	}
	if err != nil {
		return false
	}

	return false
}

func IsPhoneNumberValid(phoneNumber string) bool {
	re := regexp.MustCompile("^628[0-9]{8,14}$")
	return re.MatchString(phoneNumber)
}
