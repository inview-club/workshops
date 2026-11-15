package controllers

import "products/internal/entities"

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func (p Product) NewProduct() *entities.Product {
	return &entities.Product{
		ID:     0,
		Name:   p.Name,
		Price:  p.Price,
		Status: entities.Available,
	}
}
