package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"http/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	user := &models.User{
		Email:            "test@example.com",
		Password:         "hashedpassword",
		Is_verified:      false,
		VerificationCode: "123456",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Email, user.Password, user.Is_verified, user.VerificationCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_verified", "verification_code"}).
		AddRow(1, email, "hashedpassword", false, "123456")

	mock.ExpectQuery("SELECT id, email, password, is_verified, verification_code FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	user := &models.User{
		Id:          1,
		Is_verified: true,
	}

	mock.ExpectExec("UPDATE users SET is_verified = \\$1 WHERE id = \\$2").
		WithArgs(user.Is_verified, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	book := &models.Book{
		Title:       "Test Book",
		Author:      "Author",
		Description: "Description",
		PublishedAt: "2023-01-01",
	}

	mock.ExpectExec("INSERT INTO books").
		WithArgs(book.Title, book.Author, book.Description, book.PublishedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateBook(book)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	bookID := 1
	rows := sqlmock.NewRows([]string{"id", "title", "author", "description", "published_at"}).
		AddRow(bookID, "Test Book", "Author", "Description", "2023-01-01")

	mock.ExpectQuery("SELECT id, title, author, description, published_at FROM books WHERE id = ?").
		WithArgs(bookID).
		WillReturnRows(rows)

	book, err := repo.GetBookByID(bookID)
	assert.NoError(t, err)
	assert.Equal(t, bookID, book.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	book := &models.Book{
		ID:          1,
		Title:       "Updated Book",
		Author:      "Updated Author",
		Description: "Updated Description",
		PublishedAt: "2023-01-01",
	}

	mock.ExpectExec("UPDATE books SET title = \\$1, author = \\$2, description = \\$3, published_at = \\$4 WHERE id = \\$5").
		WithArgs(book.Title, book.Author, book.Description, book.PublishedAt, book.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateBook(book)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestDeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	bookID := 1

	mock.ExpectExec("DELETE FROM books WHERE id = ?").
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteBook(bookID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
