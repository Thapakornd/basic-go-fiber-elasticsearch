package services

import (
	"encoding/json"
	"errors"

	"example.com/m/pkg/lib"
	"example.com/m/pkg/models"
	"example.com/m/pkg/repository"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

type ProductService struct {
	db          *gorm.DB
	productRepo *repository.ProductRepository
	elastic     *lib.ElasticsearchUtil
}

func NewProductService(db *gorm.DB, es *lib.ElasticsearchUtil) *ProductService {
	productrepo := repository.NewProductRepository(db)

	return &ProductService{
		db,
		productrepo,
		es,
	}
}

func (ps *ProductService) GenerateFakeData() error {
	gofakeit.Seed(0)

	var products []*models.ProductSchema
	for i := 0; i < 1500; i++ {
		product := &models.ProductSchema{
			ID:          gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Category:    gofakeit.ProductCategory(),
			Description: gofakeit.ProductDescription(),
			Price:       gofakeit.Uint16(),
		}
		products = append(products, product)
	}
	if err := ps.productRepo.MultipleCreate(products); err != nil {
		return err
	}

	var items []models.Item
	for _, p := range products {
		items = append(items, p)
	}
	if err := ps.elastic.BlukIndexDocument(items, "products"); err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) SearchProducts(query string) ([]*models.Product, error) {
	docs, err := ps.elastic.SearchDocuments("products", query, 10)
	if err != nil {
		return nil, err
	}

	var result = []*models.Product{}
	for _, doc := range docs {
		m, ok := doc.(map[string]interface{})
		if !ok {
			return nil, errors.New("error: can't assert the doc type from elasticsearch")
		}

		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		var product models.Product
		if err := json.Unmarshal(b, &product); err != nil {
			return nil, err
		}
		result = append(result, &product)
	}

	return result, nil
}

func (ps *ProductService) InsertFakeProduct() (*models.Product, error) {
	gofakeit.Seed(0)

	price := gofakeit.Uint16()
	product := &models.Product{
		ID:          gofakeit.UUID(),
		Name:        gofakeit.ProductName(),
		Category:    gofakeit.ProductCategory(),
		Description: gofakeit.ProductDescription(),
		Price:       &price,
	}
	if err := ps.productRepo.Create(product); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	if _, err := ps.elastic.IndexDocument(&data, "products", product.ID); err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductService) UpdateProduct(id string, p *models.UpdateProductRequest) error {
	if err := ps.productRepo.Update(id, p); err != nil {
		return err
	}

	var productESReq = map[string]interface{}{}
	if p.Name != nil {
		productESReq["name"] = p.Name
	}
	if p.Category != nil {
		productESReq["category"] = p.Category
	}
	if p.Description != nil {
		productESReq["description"] = p.Description
	}
	if p.Price != nil {
		productESReq["price"] = p.Price
	}

	if _, err := ps.elastic.UpdateDocumentById("products", id, productESReq); err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) SoftDeleteProduct(id string) error {
	if err := ps.productRepo.SoftDelete(id); err != nil {
		return err
	}

	if _, err := ps.elastic.DeleteDocumentById("products", id); err != nil {
		return err
	}

	return nil
}
