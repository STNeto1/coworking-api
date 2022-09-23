package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func Enum(
	fl validator.FieldLevel,
) bool {
	enumString := fl.Param()                    // get string male_female
	value := fl.Field().String()                // the actual field
	enumSlice := strings.Split(enumString, "_") // convert to slice
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}
