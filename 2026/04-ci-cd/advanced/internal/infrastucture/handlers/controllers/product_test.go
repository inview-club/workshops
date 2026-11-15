package controllers

import (
	"testing"

	"products/internal/entities"
)

func TestNewProduct(t *testing.T) {
	ctrl := Product{Name: "apple", Price: "10"}

	got := ctrl.NewProduct()
	want := &entities.Product{ID: 0, Name: "apple", Price: "10", Status: entities.Available}

	if *got != *want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
