package usecase

import (
	"context"
	"errors"
	"fmt"
	"library_latihan/model"
)

var (
	ErrBookTitleExists   = errors.New("book title already exists")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrBookNotFound      = errors.New("book not found")
	ErrNoStock           = errors.New("no stock left")
	ErrBorrowingNotFound = errors.New("borrowing not found")
	ErrAlreadyReturned   = errors.New("book already returned")
)

type BookRepository interface {
	GetBooks(ctx context.Context, limit int, offset int, title string) ([]model.BookWithAuthorRow, error)
	IsBookTitleUnique(ctx context.Context, title string) (bool, error)
	IsAuthorExists(ctx context.Context, id int) (bool, error)
	CreateBook(ctx context.Context, book model.BookEntity) (int, error)
	BorrowBook(ctx context.Context, userID int, bookID int) error
	GetBookByID(ctx context.Context, id int) (*model.BookEntity, error)
	GetBorrowingByID(ctx context.Context, id int) (*model.BorrowingEntity, error)
	ReturnBook(ctx context.Context, borrowingID int, bookID int) error
}

type BookUsecase struct {
	repo BookRepository
}

func NewBookUsecase(repo BookRepository) *BookUsecase {
	return &BookUsecase{
		repo: repo,
	}
}

func (u *BookUsecase) GetBooks(ctx context.Context, limit int, offset int, title string) ([]model.BookDTO, error) {

	rows, err := u.repo.GetBooks(ctx, limit, offset, title)
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

func (u *BookUsecase) BorrowBook(ctx context.Context, req model.BorrowBookRequest) error {

	book, err := u.repo.GetBookByID(ctx, req.BookID)
	if err != nil {
		return fmt.Errorf("usecase get book failed: %w", err)
	}

	if book == nil {
		return ErrBookNotFound
	}

	if book.Quantity == 0 {
		return ErrNoStock
	}

	err = u.repo.BorrowBook(ctx, req.UserID, req.BookID)
	if err != nil {
		return fmt.Errorf("usecase borrow book failed: %w", err)
	}

	return nil
}

func (u *BookUsecase) ReturnBook(ctx context.Context, borrowingID int) error {

	borrowing, err := u.repo.GetBorrowingByID(ctx, borrowingID)
	if err != nil {
		return fmt.Errorf("usecase get borrowing failed: %w", err)
	}

	if borrowing == nil {
		return ErrBorrowingNotFound
	}

	if borrowing.ReturnedAt != nil {
		return ErrAlreadyReturned
	}

	err = u.repo.ReturnBook(ctx, borrowingID, borrowing.BookID)
	if err != nil {
		return fmt.Errorf("usecase return book failed: %w", err)
	}

	return nil
}
