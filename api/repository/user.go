package repository

import (
	"codis/models"
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	postgresDatabaseService *PostgresDatabaseService
}

func NewUserRepository(injector do.Injector) (*UserRepository, error) {
	u := UserRepository{
		postgresDatabaseService: do.MustInvoke[*PostgresDatabaseService](injector),
	}

	return &u, nil
}

func (u UserRepository) Create(username string, email string, password string) (models.User, error) {
	user, err := gorm.G[models.User](u.postgresDatabaseService.db).Raw(
		"INSERT INTO user (username, email, password) VALUES (?, ?, ?)",
		username,
		email,
		password,
	).First(context.Background())

	return user, err
}
