package mongodb

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "ecommerce/services/product-service/internal/domain"
)

type ProductRepository struct {
    collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
    return &ProductRepository{
        collection: db.Collection("products"),
    }
}

func (r *ProductRepository) Create(product *domain.Product) error {
    product.ID = primitive.NewObjectID()
    product.CreatedAt = time.Now()
    product.UpdatedAt = time.Now()
    product.Status = "active"

    _, err := r.collection.InsertOne(context.Background(), product)
    return err
}

func (r *ProductRepository) GetByID(id string) (*domain.Product, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var product domain.Product
    filter := bson.M{"_id": objectID, "status": bson.M{"$ne": "deleted"}}
    err = r.collection.FindOne(context.Background(), filter).Decode(&product)
    if err != nil {
        return nil, err
    }

    return &product, nil
}

func (r *ProductRepository) GetAll(limit, offset int, category, status string) ([]*domain.Product, error) {
    filter := bson.M{"status": bson.M{"$ne": "deleted"}}
    
    if category != "" {
        filter["category"] = category
    }
    if status != "" {
        filter["status"] = status
    }

    options := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
    cursor, err := r.collection.Find(context.Background(), filter, options)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var products []*domain.Product
    for cursor.Next(context.Background()) {
        var product domain.Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, &product)
    }

    return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
    product.UpdatedAt = time.Now()
    filter := bson.M{"_id": product.ID}
    update := bson.M{"$set": product}
    _, err := r.collection.UpdateOne(context.Background(), filter, update)
    return err
}

func (r *ProductRepository) Delete(id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{"status": "deleted", "updated_at": time.Now()}}
    _, err = r.collection.UpdateOne(context.Background(), filter, update)
    return err
}

func (r *ProductRepository) Search(query string, limit, offset int) ([]*domain.Product, error) {
    filter := bson.M{
        "$and": []bson.M{
            {"status": bson.M{"$ne": "deleted"}},
            {"$or": []bson.M{
                {"name": bson.M{"$regex": query, "$options": "i"}},
                {"description": bson.M{"$regex": query, "$options": "i"}},
                {"tags": bson.M{"$in": []string{query}}},
            }},
        },
    }

    options := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
    cursor, err := r.collection.Find(context.Background(), filter, options)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var products []*domain.Product
    for cursor.Next(context.Background()) {
        var product domain.Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, &product)
    }

    return products, nil
}
