package api

import (
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/zagvozdeen/ola/internal/logger"
)

var phoneRegexp = regexp.MustCompile(`^\+7 \(\d{3}\) \d{3}-\d{2}-\d{2}$`)

func newValidator(log *logger.Logger) *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.RegisterValidation("ru_phone", func(fl validator.FieldLevel) bool {
		return phoneRegexp.MatchString(fl.Field().String())
	})
	if err != nil {
		log.Error("Fatal error: failed to register custom validation tag", err)
		os.Exit(1)
	}
	return v
}
