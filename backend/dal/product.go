package dal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateProduct creates new product
func CreateProduct(p *Product, searchID primitive.ObjectID) (*Product, error) {
	c := db.Collection(ColProduct)
	ctx := context.Background()

	_, err := c.DeleteMany(ctx, bson.M{"asin": p.Asin})
	if err != nil {
		return nil, err
	}
	p.SearchID = searchID
	_, err = c.InsertOne(ctx, p)
	return p, err
}

// DeleteProduct remove one product data
func DeleteProduct(id primitive.ObjectID) error {
	c := db.Collection(ColProduct)
	ctx := context.Background()

	_, err := c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListProduct get all products
func ListProduct(searchID primitive.ObjectID) ([]Product, error) {
	ps := []Product{}
	c := db.Collection(ColProduct)
	ctx := context.Background()

	cur, err := c.Find(ctx, bson.M{"search_id": searchID})
	if err != nil {
		return ps, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &ps)
	return ps, err
}
