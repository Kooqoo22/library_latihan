package model

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