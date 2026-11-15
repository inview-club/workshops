package repository

import (
	"products/internal/entities"
)

type driver interface {
	QueryAll() ([]entities.Product, error)
	QueryByID(id string) (entities.Product, error)
	Insert(product entities.Product) (string, error)
	Ping() error
}

type database struct {
	d driver
}

func New(dbHandler driver) *database {
	return &database{
		d: dbHandler,
	}
}

func (db *database) QueryProducts() ([]entities.Product, error) {
	return db.d.QueryAll()
}

func (db *database) QueryProductByID(id string) (entities.Product, error) {
	product, err := db.d.QueryByID(id)
	if product.ID == 0 {
		product.ID = -1
	}
	return product, err
}

func (db *database) InsertProduct(product entities.Product) (string, error) {
	return db.d.Insert(product)
}

func (db *database) Ping() error {
	return db.d.Ping()
}
