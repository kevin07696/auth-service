package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

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
	if !r.db.Migrator().HasTable(&Login{}) {
		if err := r.db.Migrator().CreateTable(&Login{}); err != nil {
			slog.Error("Failed to create Login", "error", err)
			os.Exit(1)
		}
	}

	r.upsertIndexes(&Login{}, "uni_logins_username", "uni_logins_email")
}

func (r LoginRepository) upsertIndexes(dst interface{}, indexNames ...string) {
	for _, indexName := range indexNames {
		if !r.db.Migrator().HasIndex(dst, indexName) {
			if err := r.db.Migrator().CreateIndex(&Login{}, indexName); err != nil {
				message := fmt.Sprintf("Failed to create unique index %s.", indexName)
				slog.Error(message, "error", err)
				os.Exit(1)
			}
		}
	}
	slog.Info("Successfully migrated indexes.", "indexNames", indexNames)
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
	login.Username = string(username)

	status = r.getLogin(ctx, &login)
	return
}

func (r LoginRepository) GetLoginByEmail(ctx context.Context, emailAddress Email) (login Login, status StatusCode) {
	login.Email = string(emailAddress)

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
