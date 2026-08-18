package passwordpolicy

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var ErrWeak = errors.New("密码必须为 8–72 位，并包含大写字母、小写字母、数字和特殊字符")

func Validate(password string) error {
	length := utf8.RuneCountInString(password)
	if length < 8 || length > 72 {
		return ErrWeak
	}
	var upper, lower, digit, special bool
	for _, character := range password {
		switch {
		case unicode.IsUpper(character):
			upper = true
		case unicode.IsLower(character):
			lower = true
		case unicode.IsDigit(character):
			digit = true
		case unicode.IsSpace(character):
			return ErrWeak
		default:
			special = true
		}
	}
	if !upper || !lower || !digit || !special {
		return ErrWeak
	}
	return nil
}
