package domain

import (
	"regexp"
	"unicode"
)

func NewUsername(input string) (Username, bool) {
	pattern := regexp.MustCompile(`^[\w\-.]{5,20}$`)
	if isMatch := pattern.MatchString(input); !isMatch {
		return Username(""), false
	}

	return Username(input), true
}

func NewEmailComponents(input string) (EmailComponents, bool) {
	pattern := regexp.MustCompile(`^([\w.]+)[+\-]?([\w.]*)@([\w]+.[\w]+[.]?[\w]*)$`)
	if isMatch := pattern.MatchString(input); !isMatch {
		return EmailComponents{}, false
	}
	matches := pattern.FindAllStringSubmatch(input, 1)

	return EmailComponents{
		LocalAddress: matches[0][1],
		SubAddress:   matches[0][2],
		Domain:       matches[0][3],
	}, true
}

func NewPassword(input string) (Password, bool) {
	if len(input) < 12 || len(input) > 128 {
		return Password(""), false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range input {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return Password(input), hasUpper && hasLower && hasNumber && hasSpecial
}
