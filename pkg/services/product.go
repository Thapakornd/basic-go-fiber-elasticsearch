package services

import (
	"elastic-search/pkg/lib"
	"elastic-search/pkg/models"
	"elastic-search/pkg/repository"
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/elastic/go-elasticsearch/v9/esapi"
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

func (ps *ProductService) SearchProducts(query string) (*esapi.Response, error) {
	data, err := ps.elastic.SearchDocuments("products", query, 10)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", data)

	return data, nil
}

func (ps *ProductService) InsertFakeProduct() (*models.Product, error) {
	gofakeit.Seed(0)

	product := &models.Product{
		ID:          gofakeit.UUID(),
		Name:        gofakeit.ProductName(),
		Category:    gofakeit.ProductCategory(),
		Description: gofakeit.ProductDescription(),
		Price:       gofakeit.Uint16(),
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
