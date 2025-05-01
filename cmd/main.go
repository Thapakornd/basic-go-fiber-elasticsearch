package main

import (
	"elastic-search/pkg/lib"
	"elastic-search/pkg/models"
	"elastic-search/router"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env did not load -> " + err.Error())
	}

	dsn := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}))
	if err != nil {
		log.Println("Warning: database did not connect -> " + err.Error())
	}
	db.AutoMigrate(&models.Product{})

	// init-elasticsearch
	elastic := lib.NewElasticsearchUtil()
	if elastic == nil {
		log.Println("Warning: elastic did not connect")
	}

	api := app.Group("/api/v1")
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"msg": "sever is running!",
		})
	})
	router.SetupProduct(api, db, elastic)

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	app.Listen(":3000")
}
