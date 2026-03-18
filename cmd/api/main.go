package main

import (
	"context"
	"eav-intentory/internal/handler"
	"eav-intentory/internal/repository/postgres"
	"eav-intentory/internal/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// di wiring ioc container yapılacak
	connString := "postgres://postgres:123456@localhost:5432/eav_db?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	categoryRepo := postgres.NewCategoryRepository(pool)
	categoryService := usecase.NewCategoryUseCase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := postgres.NewProductRepository(pool)
	productService := usecase.NewProductUseCase(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /category", categoryHandler.Create)
	mux.HandleFunc("POST /product", productHandler.CreateProduct)
	server := http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Server Up")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
