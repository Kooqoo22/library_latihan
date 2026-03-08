package model

type AuthorEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type AuthorDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}