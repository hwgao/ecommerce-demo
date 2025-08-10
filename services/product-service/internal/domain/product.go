package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string            `json:"name" bson:"name"`
    Description string            `json:"description" bson:"description"`
    Price       float64           `json:"price" bson:"price"`
    Currency    string            `json:"currency" bson:"currency"`
    Category    string            `json:"category" bson:"category"`
    Images      []string          `json:"images" bson:"images"`
    Tags        []string          `json:"tags" bson:"tags"`
    Status      string            `json:"status" bson:"status"`
    CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

type ProductRepository interface {
    Create(product *Product) error
    GetByID(id string) (*Product, error)
    GetAll(limit, offset int, category, status string) ([]*Product, error)
    Update(product *Product) error
    Delete(id string) error
    Search(query string, limit, offset int) ([]*Product, error)
}

type ProductService interface {
    CreateProduct(name, description, category string, price float64, currency string, images, tags []string) (*Product, error)
    GetProduct(id string) (*Product, error)
    GetProducts(limit, offset int, category, status string) ([]*Product, error)
    UpdateProduct(id, name, description, category string, price float64, images, tags []string) error
    DeleteProduct(id string) error
    SearchProducts(query string, limit, offset int) ([]*Product, error)
}
