package main

import (
	"context"
	catalogHandler "eav-intentory/internal/catalog/handler"
	catalogRepository "eav-intentory/internal/catalog/repository/postgres"
	catalogUseCase "eav-intentory/internal/catalog/usecase"
	inventory_handler "eav-intentory/internal/inventory/handler"
	inventory_repository "eav-intentory/internal/inventory/repository"
	inventory_usecase "eav-intentory/internal/inventory/usecase"
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
	categoryRepo := catalogRepository.NewCategoryRepository(pool)
	categoryService := catalogUseCase.NewCategoryUseCase(categoryRepo)
	categoryHandler := catalogHandler.NewCategoryHandler(categoryService)

	productRepo := catalogRepository.NewProductRepository(pool)
	productService := catalogUseCase.NewProductUseCase(productRepo, categoryRepo)
	productHandler := catalogHandler.NewProductHandler(productService)

	attributeRepo := catalogRepository.NewAttributeRepository(pool)
	attributeService := catalogUseCase.NewAttributeService(attributeRepo)
	attributeHandler := catalogHandler.NewAttributeHandler(attributeService)

	// todo : importlarda hala problem var warehouse katmanlarını alıp di bağlantısını yapacağım fakat importlar sinirimi bozdu
	warehouseRepo := inventory_repository.NewWarehouseRepository(pool)
	warehouseService := inventory_usecase.NewWarehouseUsecase(warehouseRepo)
	warehouseHandler := inventory_handler.NewWarehouseHandler(warehouseService)
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
	mux.HandleFunc("PUT /category/{categoryId}/attribute/{attributeId}", categoryHandler.UpdateAttributeFromCategory)
	mux.HandleFunc("DELETE /category/{categoryId}/attribute/{attributeId}", categoryHandler.RemoveAttributeFromCategory)

	//products
	mux.HandleFunc("POST /product", productHandler.CreateProduct)
	mux.HandleFunc("GET /product/{id}", productHandler.GetById)
	mux.HandleFunc("GET /product/", productHandler.GetProducts)

	mux.HandleFunc("POST /warehouse", warehouseHandler.CreateWarehouse)

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
