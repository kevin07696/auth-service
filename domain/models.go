package domain

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// A validated username string.
type Username string

// A validated email string.
type Email string

// A validated password string.
type Password []byte

// Represents the "logins" table in the GORM database.
// Unique Indexes: username, email.
// Index: deleted_at.
type Login struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex:uni_logins_username; not null"`
	Email          string `gorm:"uniqueIndex:uni_logins_email; not null"`
	HashedPassword string `gorm:"not null"`
}

// Hashes a password with a hasher adapter
func (p Password) HashPassword(hasher Hasher) Password {
	return Password(hasher.HashPassword(p))
}

// Verifies a password  with a hasher adapter
func (p Password) VerifyPassword(hasher Hasher, hashedPassword Password) bool {
	return hasher.VerifyPassword(hashedPassword, p)
}

type EmailComponents struct {
	LocalAddress string
	SubAddress   string
	Domain       string
}

// Returns an email address without a subaddress.
// Standard is used for verifying unique email addresses at login.
func (e EmailComponents) ToStandardString() Email {
	return Email(fmt.Sprintf("%s@%s", e.LocalAddress, e.Domain))
}

// Returns a full email address with a subaddress.
// Full is used for sending emails to user
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
