package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Username string
type Email string
type Password []byte

type Base struct {
	ID        string     `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (u *Base) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()

	return
}

type Login struct {
	Base
	Username             Username `gorm:"uniqueIndex;not null"`
	StandardEmailAddress Email    `gorm:"uniqueIndex;not null"`
	HashedPassword       Password `gorm:"not null"`
}

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
