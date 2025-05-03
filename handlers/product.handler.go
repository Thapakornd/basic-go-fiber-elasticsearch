package handlers

import (
	"example.com/m/pkg/lib"
	"example.com/m/pkg/models"
	"example.com/m/pkg/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db             *gorm.DB
	productService *services.ProductService
}

func NewProductHandler(db *gorm.DB, es *lib.ElasticsearchUtil) *ProductHandler {
	productsvc := services.NewProductService(db, es)

	return &ProductHandler{
		db,
		productsvc,
	}
}

func (ph *ProductHandler) GenerateFakeData(c *fiber.Ctx) error {
	if err := ph.productService.GenerateFakeData(); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"msg":    err.Error(),
			"code":   10001,
			"status": fiber.ErrBadRequest,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"msg": "data generated!",
	})
}

func (ph *ProductHandler) Create(c *fiber.Ctx) error {
	product, err := ph.productService.InsertFakeProduct()
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"msg":    err.Error(),
			"code":   20001,
			"status": fiber.ErrBadRequest.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (ph *ProductHandler) Search(c *fiber.Ctx) error {
	var req models.ProductSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"msg":    err.Error(),
			"code":   30001,
			"status": fiber.ErrBadRequest.Code,
		})
	}

	res, err := ph.productService.SearchProducts(req.Query)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"msg":    err.Error(),
			"code":   30002,
			"status": fiber.ErrBadRequest.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
