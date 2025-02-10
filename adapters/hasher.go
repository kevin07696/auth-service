package adapters

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type bcryptAdapter struct {
	cost int
}

func NewBcryptAdapter(cost int) bcryptAdapter {
	return bcryptAdapter{cost: cost}
}

func (b bcryptAdapter) HashPassword(password []byte) string {
	password, err := bcrypt.GenerateFromPassword(password, b.cost)
	if err != nil {
		log.Fatalf("Failed to hash password: %s. Error: %v", string(password), err)
	}
	return string(password)
}

func (b bcryptAdapter) VerifyPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	return err != nil
}
