package main

import (
	"log"
	"time"

	"library_latihan/config"
	"library_latihan/handler"
	"library_latihan/middleware"
	"library_latihan/repository"
	"library_latihan/usecase"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.New()

	router.Use(middleware.Logging())
	router.Use(middleware.Timeout(5 * time.Second))
	router.Use(gin.Recovery())

	// dependency injection
	bookRepo := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepo)
	bookHandler := handler.NewBookHandler(bookUsecase)

	// API group
	api := router.Group("/api")

	// version group
	v1 := api.Group("/v1")

	books := v1.Group("/books")
	{
		books.GET("", bookHandler.GetBooks)

		admin := books.Group("")
		admin.Use(middleware.Auth(), middleware.AdminOnly())
		{
			admin.POST("", bookHandler.CreateBook)
		}
	}
	log.Println("server running at :8080")

	router.Run(":8080")
}
