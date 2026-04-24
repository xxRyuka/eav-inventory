package main

import (
	"context"
	catalog3 "eav-intentory/internal/handler/catalog"
	catalog2 "eav-intentory/internal/repository/postgres/catalog"
	"eav-intentory/internal/usecase/catalog"
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
	categoryRepo := catalog2.NewCategoryRepository(pool)
	categoryService := catalog.NewCategoryUseCase(categoryRepo)
	categoryHandler := catalog3.NewCategoryHandler(categoryService)

	productRepo := catalog2.NewProductRepository(pool)
	productService := catalog.NewProductUseCase(productRepo, categoryRepo)
	productHandler := catalog3.NewProductHandler(productService)

	attributeRepo := catalog2.NewAttributeRepository(pool)
	attributeService := catalog.NewAttributeService(attributeRepo)
	attributeHandler := catalog3.NewAttributeHandler(attributeService)

	mux := http.NewServeMux()

	//attributes
	mux.HandleFunc("POST /attribute", attributeHandler.CreateAttribute)
	mux.HandleFunc("GET /attribute/{id}", attributeHandler.GetAttributeByID)
	mux.HandleFunc("GET /attribute/", attributeHandler.GetAttributes)
	mux.HandleFunc("POST /attribute/{id}", attributeHandler.DeleteAttribute)
	mux.HandleFunc("PUT /attribute/", attributeHandler.UpdateAttribute)

	//categories
	mux.HandleFunc("POST /category", categoryHandler.Create)
	mux.HandleFunc("GET /category/{id}", categoryHandler.GetCategoryById)
	mux.HandleFunc("GET /category", categoryHandler.GetCategories)
	mux.HandleFunc("PUT /category/{id}", categoryHandler.UpdateBaseCategory)
	mux.HandleFunc("GET /categories/detailed", categoryHandler.GetCategoriesWithAttributes)
	mux.HandleFunc("POST /category/assign-attribute", categoryHandler.AssignAttributeToCategory)

	//products
	mux.HandleFunc("POST /product", productHandler.CreateProduct)
	mux.HandleFunc("GET /product/{id}", productHandler.GetById)
	mux.HandleFunc("GET /product/", productHandler.GetProducts)

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
