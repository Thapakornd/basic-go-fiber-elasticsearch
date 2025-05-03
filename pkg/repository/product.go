package repository

import (
	"context"
	"time"

	"example.com/m/pkg/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

func (pr *ProductRepository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 30*time.Second)
}

func (pr *ProductRepository) MultipleCreate(products []*models.ProductSchema) error {
	ctx, cancel := pr.withTimeout(context.Background())
	defer cancel()

	if err := pr.db.WithContext(ctx).Table("products").Create(products).Error; err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) Create(product *models.Product) error {
	ctx, cancel := pr.withTimeout(context.Background())
	defer cancel()

	if err := pr.db.WithContext(ctx).Table("products").Create(product).Error; err != nil {
		return err
	}

	return nil
}
