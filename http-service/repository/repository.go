package repository

import (
	"database/sql"
	"http/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password, is_verified) VALUES ($1, $2, $3)", user.Email, user.Password, user.Is_verified)
	return err
}
