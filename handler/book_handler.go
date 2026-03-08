package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"library_latihan/model"
	"library_latihan/usecase"

	"errors"

	"github.com/gin-gonic/gin"
)

type BookUsecase interface {
	GetBooks(ctx context.Context, limit int, offset int, title string) ([]model.BookDTO, error)
	CreateBook(ctx context.Context, req model.CreateBookRequest) (int, error)
	BorrowBook(ctx context.Context, req model.BorrowBookRequest) error
	ReturnBook(ctx context.Context, borrowingID int) error
}

type BookHandler struct {
	usecase BookUsecase
}

func NewBookHandler(usecase BookUsecase) *BookHandler {
	return &BookHandler{
		usecase: usecase,
	}
}

func (h *BookHandler) GetBooks(c *gin.Context) {

	ctx := c.Request.Context()

	limitStr, limitExists := c.GetQuery("limit")
	offsetStr, offsetExists := c.GetQuery("offset")
	title := c.Query("title")

	var limit int
	var offset int
	var err error

	if limitExists {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
	}

	if offsetExists {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
	}

	books, err := h.usecase.GetBooks(ctx, limit, offset, title)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get books",
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) CreateBook(c *gin.Context) {

	ctx := c.Request.Context()

	var req model.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	id, err := h.usecase.CreateBook(ctx, req)
	if err != nil {

		// LOG FULL ERROR CHAIN
		log.Printf("CreateBook error: %+v\n", err)

		if errors.Is(err, usecase.ErrBookTitleExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, usecase.ErrAuthorNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create book",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *BookHandler) BorrowBook(c *gin.Context) {

	ctx := c.Request.Context()

	var req model.BorrowBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.usecase.BorrowBook(ctx, req)
	if err != nil {

		if errors.Is(err, usecase.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, usecase.ErrNoStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "book borrowed successfully",
	})
}

func (h *BookHandler) ReturnBook(c *gin.Context) {

	ctx := c.Request.Context()

	var req model.ReturnBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.usecase.ReturnBook(ctx, req.BorrowingID)
	if err != nil {

		if errors.Is(err, usecase.ErrBorrowingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, usecase.ErrAlreadyReturned) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to return book",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "book returned successfully",
	})
}
