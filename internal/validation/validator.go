package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		re := regexp.MustCompile(`^[0-9+\-()\s]{7,20}$`)
		return re.MatchString(s)
	})
	_ = validate.RegisterValidation("cvv", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		re := regexp.MustCompile(`^\d{3,4}$`)
		return re.MatchString(s)
	})
	_ = validate.RegisterValidation("month", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		re := regexp.MustCompile(`^(?:0?[1-9]|1[0-2])$`)
		return re.MatchString(s)
	})
	_ = validate.RegisterValidation("card", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		for _, r := range s {
			if !unicode.IsDigit(r) {
				return false
			}
		}
		if len(s) < 13 || len(s) > 19 {
			return false
		}
		sum := 0
		dbl := false
		for i := len(s) - 1; i >= 0; i-- {
			d := int(s[i] - '0')
			if dbl {
				d *= 2
				if d > 9 {
					d -= 9
				}
			}
			sum += d
			dbl = !dbl
		}
		return sum%10 == 0
	})
}

func Get() *validator.Validate {
	return validate
}

func ValidateStruct(v any) error {
	return validate.Struct(v)
}
