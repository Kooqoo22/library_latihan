package handler

import (
	"context"
	"log"
	"net/http"

	"library_latihan/model"
	"library_latihan/usecase"

	"errors"

	"github.com/gin-gonic/gin"
)

type BookUsecase interface {
	GetBooks(ctx context.Context) ([]model.BookDTO, error)
	CreateBook(ctx context.Context, req model.CreateBookRequest) (int, error)
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

	books, err := h.usecase.GetBooks(ctx)
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
