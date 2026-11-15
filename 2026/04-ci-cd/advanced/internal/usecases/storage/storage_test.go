package storage

import (
	"errors"
	"testing"

	"products/internal/entities"
)

type fakeRepo struct {
	products []entities.Product
	product  entities.Product
	id       string
	err      error
	pingErr  error
	inserted entities.Product
}

func (f *fakeRepo) QueryProducts() ([]entities.Product, error) {
	return f.products, f.err
}

func (f *fakeRepo) QueryProductByID(id string) (entities.Product, error) {
	return f.product, f.err
}

func (f *fakeRepo) InsertProduct(product entities.Product) (string, error) {
	f.inserted = product
	return f.id, f.err
}

func (f *fakeRepo) Ping() error {
	return f.pingErr
}

func TestGetProducts(t *testing.T) {
	want := []entities.Product{{ID: 1, Name: "apple", Price: "10", Status: entities.Available}}
	app := New(&fakeRepo{products: want})

	got, err := app.GetProducts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "apple" {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetProductsError(t *testing.T) {
	repoErr := errors.New("db error")
	app := New(&fakeRepo{err: repoErr})

	if _, err := app.GetProducts(); !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}

func TestGetProductByID(t *testing.T) {
	want := entities.Product{ID: 7, Name: "pear", Price: "5"}
	app := New(&fakeRepo{product: want})

	got, err := app.GetProductByID("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestAddProduct(t *testing.T) {
	repo := &fakeRepo{id: "42"}
	app := New(repo)

	product := entities.Product{Name: "plum", Price: "3", Status: entities.Available}
	id, err := app.AddProduct(product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "42" {
		t.Errorf("got id %q, want %q", id, "42")
	}
	if repo.inserted != product {
		t.Errorf("inserted %+v, want %+v", repo.inserted, product)
	}
}

func TestHealth(t *testing.T) {
	if err := New(&fakeRepo{}).Health(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	pingErr := errors.New("db unreachable")
	if err := New(&fakeRepo{pingErr: pingErr}).Health(); !errors.Is(err, pingErr) {
		t.Errorf("got %v, want %v", err, pingErr)
	}
}
