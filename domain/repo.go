package domain

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type LoginRepository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) (repo LoginRepository) {
	repo.db = db
	return
}

func (r LoginRepository) errorHandler(err error) StatusCode {
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return StatusDuplicateKey
	case errors.Is(err, gorm.ErrRecordNotFound):
		return StatusNotFound
	default:
		return StatusInternal
	}
}

func (r LoginRepository) Migrate() {
	if err := r.db.AutoMigrate(&Login{}); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}
}

func (r LoginRepository) CreateLogin(ctx context.Context, login *Login) (status StatusCode) {
	result := r.db.WithContext(ctx).Create(login)
	if result.Error != nil {
		status = r.errorHandler(result.Error)
		return
	}

	return
}

func (r LoginRepository) GetLoginByUsername(ctx context.Context, username Username) (login Login, status StatusCode) {
	login.Username = username

	status = r.getLogin(ctx, &login)
	return
}

func (r LoginRepository) GetLoginByEmail(ctx context.Context, emailAddress Email) (login Login, status StatusCode) {
	login.StandardEmailAddress = emailAddress

	status = r.getLogin(ctx, &login)
	return
}

func (r LoginRepository) getLogin(ctx context.Context, login *Login) (status StatusCode) {
	result := r.db.WithContext(ctx).Take(login)
	if result.Error != nil {
		status = r.errorHandler(result.Error)
		return
	}

	return
}
