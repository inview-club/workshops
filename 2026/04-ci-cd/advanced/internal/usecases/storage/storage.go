package storage

import "products/internal/entities"

type repository interface {
	QueryProducts() ([]entities.Product, error)
	QueryProductByID(id string) (entities.Product, error)
	InsertProduct(product entities.Product) (string, error)
	Ping() error
}

type Controller interface {
	GetProducts() ([]entities.Product, error)
	GetProductByID(id string) (entities.Product, error)
	AddProduct(product entities.Product) (string, error)
	Health() error
}

type application struct {
	repo repository
}

func New(repo repository) *application {
	return &application{
		repo: repo,
	}
}

func (app *application) AddProduct(product entities.Product) (string, error) {
	return app.repo.InsertProduct(product)
}

func (app *application) GetProducts() ([]entities.Product, error) {
	return app.repo.QueryProducts()
}

func (app *application) GetProductByID(id string) (entities.Product, error) {
	return app.repo.QueryProductByID(id)
}

func (app *application) Health() error {
	return app.repo.Ping()
}
