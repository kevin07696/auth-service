package domain

type Hasher interface {
	HashPassword(password []byte) string
	VerifyPassword(hashedPassword, password []byte) bool
}
