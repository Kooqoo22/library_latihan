package usecase

import (
	"context"
	"errors"
	"fmt"
	"library_latihan/model"
)

var (
	ErrBookTitleExists = errors.New("book title already exists")
	ErrAuthorNotFound  = errors.New("author not found")
)

type BookRepository interface {
	GetBooks(ctx context.Context) ([]model.BookWithAuthorRow, error)
	IsBookTitleUnique(ctx context.Context, title string) (bool, error)
	IsAuthorExists(ctx context.Context, id int) (bool, error)
	CreateBook(ctx context.Context, book model.BookEntity) (int, error)
}

type BookUsecase struct {
	repo BookRepository
}

func NewBookUsecase(repo BookRepository) *BookUsecase {
	return &BookUsecase{
		repo: repo,
	}
}

func (u *BookUsecase) GetBooks(ctx context.Context) ([]model.BookDTO, error) {

	rows, err := u.repo.GetBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase GetBooks failed: %w", err)
	}

	var books []model.BookDTO

	for _, r := range rows {

		book := model.BookDTO{
			ID:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			Quantity:    r.Quantity,
			Cover:       r.Cover,
			Author: model.AuthorDTO{
				ID:   r.AuthorID,
				Name: r.AuthorName,
			},
		}

		books = append(books, book)
	}

	return books, nil
}

func (u *BookUsecase) CreateBook(ctx context.Context, req model.CreateBookRequest) (int, error) {

	unique, err := u.repo.IsBookTitleUnique(ctx, req.Title)
	if err != nil {
		return 0, fmt.Errorf("usecase check book title uniqueness failed: %w", err)
	}

	if !unique {
		return 0, ErrBookTitleExists
	}

	exists, err := u.repo.IsAuthorExists(ctx, req.AuthorID)
	if err != nil {
		return 0, fmt.Errorf("usecase check author exists failed: %w", err)
	}

	if !exists {
		return 0, ErrAuthorNotFound
	}

	book := model.BookEntity{
		Title:    req.Title,
		AuthorID: req.AuthorID,
	}

	if req.Description != nil {
		book.Description = *req.Description
	}

	if req.Quantity != nil {
		book.Quantity = *req.Quantity
	}

	if req.Cover != nil {
		book.Cover = *req.Cover
	}

	id, err := u.repo.CreateBook(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("usecase create book failed: %w", err)
	}

	return id, nil
}
