package dal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateStandard creates new standard product
func CreateStandard(p Standard) error {
	c := db.Collection(ColStandard)
	ctx := context.Background()
	_, err := c.InsertOne(ctx, p)
	return err
}

// UpdateStandard updates existing standard product
func UpdateStandard(p Standard) error {
	c := db.Collection(ColStandard)
	ctx := context.Background()

	_, err := c.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$set": bson.M{
		"price": p.Price,
		"date":  p.Date,
	}})
	return err
}

// GetStandardByAsin get a standard product by asin
func GetStandardByAsin(asin string) (*Standard, error) {
	p := &Standard{}
	c := db.Collection(ColStandard)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{"asin": asin}).Decode(&p)
	return p, err
}
