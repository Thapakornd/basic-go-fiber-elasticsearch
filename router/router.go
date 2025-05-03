package router

import (
	"example.com/m/handlers"
	"example.com/m/pkg/lib"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupProduct(api fiber.Router, db *gorm.DB, es *lib.ElasticsearchUtil) {
	productHandler := handlers.NewProductHandler(db, es)

	api.Group("/products").
		Get("/init", productHandler.GenerateFakeData).
		Post("/create", productHandler.Create).
		Post("/search", productHandler.Search)
}
