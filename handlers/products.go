// Package classification of Product API
//
// Documentation of Product API
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sheetalrprasad/productapi/data"
)

// A list of products returns in the response
//sawgger:response productResponse
type productsResponse struct {
	// All produtcs in the system
	// in:body
	Body []data.Product
}

// Products in a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200:productsResponse

// GetProducts returns a list of all products in the system
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Failed to marshal json.", http.StatusInternalServerError)
	}
}

// AddProducts creates a new product to the system
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling the POST")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201:noContentResponse
//  404:errorResponse
//  501:errorResponse

// Delete handles DELETE requests and removes items from the database
// DeleteProducts creates a new product to the system
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling the DELETE")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Cannot Caonvert ID to int", http.StatusBadRequest)
	}
	err = data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product data not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product data not found", http.StatusInternalServerError)
		return
	}
}

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201:noContentResponse
//  404:errorResponse
//  422:errorValidation

// UpdateProduct modifies an existing product in the system
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Cannot Caonvert ID to int", http.StatusBadRequest)
	}
	p.l.Println("Handling the PUT ", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product data not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product data not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Cannot unmarshal", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}
		p.l.Printf("prod %#v", prod) //shows field + value
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
