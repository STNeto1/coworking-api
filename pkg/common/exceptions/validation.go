package exceptions

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "datetime":
		return "Invalid date"
	case "min":
		return "This field must be at least " + fe.Param() + " characters"
	case "max":
		return "This field must be at most " + fe.Param() + " characters"
	case "Enum":
		replacer := *strings.NewReplacer("_", ", ")
		return fe.Field() + " must be one of " + replacer.Replace(fe.Param())

	}

	return fe.Error() // default error
}
