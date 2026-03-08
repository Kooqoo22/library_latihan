package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *BookRepository) GetBooks(ctx context.Context, limit int, offset int, title string) ([]model.BookWithAuthorRow, error) {

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

	args := []interface{}{}
	argPos := 1

	if title != "" {
		query += fmt.Sprintf(" AND b.title ILIKE $%d", argPos)
		args = append(args, "%"+title+"%")
		argPos++
	}

	query += " ORDER BY b.id"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
		argPos++
	}

	var rows []model.BookWithAuthorRow

	err := r.DB.SelectContext(ctx, &rows, query, args...)
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

func (r *BookRepository) BorrowBook(ctx context.Context, userID int, bookID int) error {

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// insert borrowing record
	insertQuery := `
		INSERT INTO borrowings (user_id, book_id, status)
		VALUES ($1,$2,'borrowed')
	`

	_, err = tx.ExecContext(ctx, insertQuery, userID, bookID)
	if err != nil {
		return err
	}

	// decrease book quantity
	updateQuery := `
	UPDATE books
	SET quantity = quantity - 1
	WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateQuery, bookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *BookRepository) GetBookByID(ctx context.Context, id int) (*model.BookEntity, error) {

	query := `
	SELECT id, title, quantity
	FROM books
	WHERE id = $1 AND deleted_at IS NULL
	`

	var book model.BookEntity

	err := r.DB.GetContext(ctx, &book, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) GetBorrowingByID(ctx context.Context, id int) (*model.BorrowingEntity, error) {

	query := `
	SELECT id, book_id, returned_at
	FROM borrowings
	WHERE id = $1
	`

	var borrowing model.BorrowingEntity

	err := r.DB.GetContext(ctx, &borrowing, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &borrowing, nil
}

func (r *BookRepository) ReturnBook(ctx context.Context, borrowingID int, bookID int) error {

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// update borrowing
	updateBorrow := `
		UPDATE borrowings
		SET status = 'returned',
		returned_at = NOW()
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateBorrow, borrowingID)
	if err != nil {
		return err
	}

	// increase quantity
	updateBook := `
		UPDATE books
		SET quantity = quantity + 1
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateBook, bookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
