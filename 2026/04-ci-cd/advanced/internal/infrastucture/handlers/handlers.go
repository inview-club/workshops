package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"products/internal/infrastucture/handlers/controllers"
	"products/internal/infrastucture/handlers/middleware"
	"products/internal/infrastucture/handlers/views"
	"products/internal/usecases/storage"
)

func getProducts(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error get products"
		products, err := app.GetProducts()

		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := make(map[string]interface{})
		response["products"] = products
		json.NewEncoder(w).Encode(response)
	})
}

func getProductByID(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error get product"
		params := mux.Vars(r)
		id := params["id"]
		product, err := app.GetProductByID(id)

		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		if product.ID == -1 {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := make(map[string]interface{})
		response["product"] = product
		json.NewEncoder(w).Encode(response)
	})
}

func addProduct(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error add product"

		ctrlProduct := new(controllers.Product)
		err := json.NewDecoder(r.Body).Decode(&ctrlProduct)
		if err != nil {
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		eProduct := ctrlProduct.NewProduct()
		id, err := app.AddProduct(*eProduct)
		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&views.ID{ID: id})
	})
}

func health(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := app.Health(); err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

func Make(r *mux.Router, app storage.Controller) {
	apiURI := "/api"
	serviceRouter := r.PathPrefix(apiURI).Subrouter()
	r.Use(middleware.MetricsMiddleware)
	r.MethodNotAllowedHandler = NotAllowedHandler()
	r.NotFoundHandler = NotFoundHandler()
	r.Handle("/health", health(app)).Methods("GET")
	serviceRouter.Handle("/products", getProducts(app)).Methods("GET")
	serviceRouter.Handle("/products/{id}", getProductByID(app)).Methods("GET")
	serviceRouter.Handle("/products", addProduct(app)).Methods("POST")
}
