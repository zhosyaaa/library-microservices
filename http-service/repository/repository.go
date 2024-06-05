package repository

import (
	"database/sql"
	"fmt"
	"http/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password, is_verified, verification_code) VALUES ($1, $2, $3, $4)", user.Email, user.Password, user.Is_verified, user.VerificationCode)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, email, password, is_verified, verification_code FROM users WHERE email = $1", email).Scan(&user.Id, &user.Email, &user.Password, &user.Is_verified, &user.VerificationCode)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec("UPDATE users SET is_verified = $1 WHERE id = $2", true, user.Id)
	return err
}
func (r *Repository) CreateBook(book *models.Book) error {
	_, err := r.db.Exec("INSERT INTO books (title, author, description, published_at) VALUES ($1, $2, $3, $4)", book.Title, book.Author, book.Description, book.PublishedAt)
	return err
}

func (r *Repository) GetBookByID(id int) (*models.Book, error) {
	book := &models.Book{}
	err := r.db.QueryRow("SELECT id, title, author, description, published_at FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.PublishedAt)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *Repository) GetBooks(filter, sort string, limit, offset int) ([]*models.Book, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT id, title, author, description, published_at FROM books"

	if filter != "" {
		query += fmt.Sprintf(" WHERE title ILIKE '%%%s%%' OR author ILIKE '%%%s%%'", filter, filter)
	}

	if sort != "" {
		query += " ORDER BY " + sort
	} else {
		query += " ORDER BY id"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err = r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.PublishedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *Repository) UpdateBook(book *models.Book) error {
	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2, description = $3, published_at = $4 WHERE id = $5", book.Title, book.Author, book.Description, book.PublishedAt, book.ID)
	return err
}

func (r *Repository) DeleteBook(id int) error {
	_, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}
