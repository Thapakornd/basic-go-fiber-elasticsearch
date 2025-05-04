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

func (pr *ProductRepository) Update(id string, product *models.UpdateProductRequest) error {
	ctx, cancel := pr.withTimeout(context.Background())
	defer cancel()

	tx := pr.db.WithContext(ctx).Table("products").Where("id = ? AND (is_delete = 0 OR is_delete IS NULL)", id)
	updates := map[string]interface{}{}

	if product.Name != nil {
		updates["name"] = *product.Name
	}
	if product.Category != nil {
		updates["category"] = *product.Category
	}
	if product.Description != nil {
		updates["description"] = *product.Description
	}
	if product.Price != nil {
		updates["price"] = *product.Price
	}

	if len(updates) > 0 {
		if err := tx.Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (pr *ProductRepository) SoftDelete(id string) error {
	ctx, cancel := pr.withTimeout(context.Background())
	defer cancel()

	if err := pr.db.WithContext(ctx).Table("products").
		Where("id = ?", id).
		Update("is_delete", 1).Error; err != nil {
		return err
	}

	return nil
}
