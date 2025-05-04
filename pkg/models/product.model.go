package models

type ProductSchema struct {
	Item
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Price       uint16 `json:"price"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       *uint16 `json:"price"`
	IsDelete    bool    `json:"products" gorm:"default:0"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Category    *string `json:"category"`
	Description *string `json:"description"`
	Price       *uint16 `json:"price"`
}

func (ps *ProductSchema) GetID() string { return ps.ID }

type ProductSearchRequest struct {
	Query string `json:"query"`
}
