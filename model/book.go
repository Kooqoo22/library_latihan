package model

import "time"

type BookEntity struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Quantity    int    `db:"quantity"`
	Cover       string `db:"cover"`
	AuthorID    int    `db:"author_id"`
}

type BookDTO struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Cover       string    `json:"cover"`
	Author      AuthorDTO `json:"author"`
}

type BookWithAuthorRow struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Quantity    int    `db:"quantity"`
	Cover       string `db:"cover"`

	AuthorID   int    `db:"author_id"`
	AuthorName string `db:"author_name"`
}

type CreateBookRequest struct {
	Title       string  `json:"title" binding:"required"`
	AuthorID    int     `json:"author_id" binding:"required"`
	Description *string `json:"description"`
	Quantity    *int    `json:"quantity"`
	Cover       *string `json:"cover"`
}

type BorrowBookRequest struct {
	UserID int `json:"user_id" binding:"required"`
	BookID int `json:"book_id" binding:"required"`
}

type ReturnBookRequest struct {
	BorrowingID int `json:"borrowing_id" binding:"required"`
}

type BorrowingEntity struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	BookID     int        `db:"book_id"`
	BorrowedAt time.Time  `db:"borrowed_at"`
	ReturnedAt *time.Time `db:"returned_at"`
}
