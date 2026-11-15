package entities

type Product struct {
	ID     int
	Name   string
	Status ProductStatus
	Price  string
}

type ProductStatus int

const (
	Available ProductStatus = iota + 1
	NotAvailable
)
