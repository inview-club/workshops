package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"products/internal/entities"
)

type fakeApp struct {
	products []entities.Product
	product  entities.Product
	id       string
	err      error
	health   error
}

func (f *fakeApp) GetProducts() ([]entities.Product, error) {
	return f.products, f.err
}

func (f *fakeApp) GetProductByID(id string) (entities.Product, error) {
	return f.product, f.err
}

func (f *fakeApp) AddProduct(product entities.Product) (string, error) {
	return f.id, f.err
}

func (f *fakeApp) Health() error {
	return f.health
}

func newTestRouter(app *fakeApp) *mux.Router {
	r := mux.NewRouter()
	Make(r, app)
	return r
}

func TestGetProducts(t *testing.T) {
	app := &fakeApp{products: []entities.Product{{ID: 1, Name: "apple", Price: "10"}}}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string][]entities.Product
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body["products"]) != 1 || body["products"][0].Name != "apple" {
		t.Errorf("unexpected products in response: %+v", body["products"])
	}
}

func TestGetProductsError(t *testing.T) {
	app := &fakeApp{err: errors.New("db error")}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestGetProductByID(t *testing.T) {
	app := &fakeApp{product: entities.Product{ID: 3, Name: "pear", Price: "5"}}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/api/products/3", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]entities.Product
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["product"].ID != 3 {
		t.Errorf("unexpected product in response: %+v", body["product"])
	}
}

func TestGetProductByIDNotFound(t *testing.T) {
	app := &fakeApp{product: entities.Product{ID: -1}}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/api/products/99", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestAddProduct(t *testing.T) {
	app := &fakeApp{id: "9"}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(`{"name":"plum","price":"3"}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["id"] != "9" {
		t.Errorf("got id %q, want %q", body["id"], "9")
	}
}

func TestAddProductBadRequest(t *testing.T) {
	app := &fakeApp{}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(`not json`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHealth(t *testing.T) {
	app := &fakeApp{}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHealthUnavailable(t *testing.T) {
	app := &fakeApp{health: errors.New("db unreachable")}
	r := newTestRouter(app)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestNotFoundHandler(t *testing.T) {
	r := newTestRouter(&fakeApp{})

	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
