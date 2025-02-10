package domain

import (
	"fmt"
	"strings"
)

type Username string
type Email string
type Password []byte

func (p Password) HashPassword(hasher Hasher) Password {
	return Password(hasher.HashPassword(p))
}

func (p Password) VerifyPassword(hasher Hasher, hashedPassword Password) bool {
	return hasher.VerifyPassword(hashedPassword, p)
}

type EmailComponents struct {
	LocalAddress string
	SubAddress   string
	Domain       string
}

func (e EmailComponents) ToStandardString() Email {
	return Email(fmt.Sprintf("%s@%s", e.LocalAddress, e.Domain))
}

func (e EmailComponents) ToFullString() Email {
	var builder strings.Builder
	builder.WriteString(e.LocalAddress)
	if e.SubAddress != "" {
		builder.WriteString("+")
		builder.WriteString(e.SubAddress)
	}
	builder.WriteString("@")
	builder.WriteString(e.Domain)
	return Email(builder.String())
}