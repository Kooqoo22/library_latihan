package repository

import (
	"context"
	"fmt"
	"library_latihan/model"

	"github.com/jmoiron/sqlx"
)

type BookRepository struct {
	DB *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{
		DB: db,
	}
}

func (r *BookRepository) GetBooks(ctx context.Context) ([]model.BookWithAuthorRow, error) {

	query := `
	SELECT
		b.id,
		b.title,
		b.description,
		b.quantity,
		b.cover,
		a.id AS author_id,
		a.name AS author_name
	FROM books b
	JOIN authors a ON b.author_id = a.id
	WHERE b.deleted_at IS NULL
	`

	var rows []model.BookWithAuthorRow

	err := r.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("repository GetBooks query failed: %w", err)
	}

	return rows, nil
}

func (r *BookRepository) IsBookTitleUnique(ctx context.Context, title string) (bool, error) {

	query := `
	SELECT COUNT(1)
	FROM books
	WHERE title = $1
	AND deleted_at IS NULL
	`

	var count int

	err := r.DB.GetContext(ctx, &count, query, title)
	if err != nil {
		return false, fmt.Errorf("repository IsBookTitleUnique query failed: %w", err)
	}

	return count == 0, nil
}

func (r *BookRepository) IsAuthorExists(ctx context.Context, id int) (bool, error) {

	query := `
	SELECT COUNT(1)
	FROM authors
	WHERE id = $1
	AND deleted_at IS NULL
	`

	var count int

	err := r.DB.GetContext(ctx, &count, query, id)
	if err != nil {
		return false, fmt.Errorf("repository IsAuthorExists query failed: %w", err)
	}

	return count > 0, nil
}

func (r *BookRepository) CreateBook(ctx context.Context, book model.BookEntity) (int, error) {

	query := `
	INSERT INTO books
	(title, description, quantity, cover, author_id)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	var id int

	err := r.DB.QueryRowContext(
		ctx,
		query,
		book.Title,
		book.Description,
		book.Quantity,
		book.Cover,
		book.AuthorID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository CreateBook insert failed: %w", err)
	}

	return id, nil
}
