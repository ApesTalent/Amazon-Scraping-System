package dal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateProxy creates new proxy string
func CreateProxy(p Proxy) error {
	c := db.Collection(ColProxy)
	ctx := context.Background()
	_, err := c.InsertOne(ctx, p)
	return err
}

// DeleteProxy remove one proxy string
func DeleteProxy(id primitive.ObjectID) error {
	c := db.Collection(ColProxy)
	ctx := context.Background()

	_, err := c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteProxy remove one proxy string
func DeleteProxies() error {
	c := db.Collection(ColProxy)
	ctx := context.Background()

	_, err := c.DeleteMany(ctx, bson.M{})
	return err
}

// ListProxy get all proxies
func ListProxy() ([]Proxy, error) {
	ps := []Proxy{}
	c := db.Collection(ColProxy)
	ctx := context.Background()

	cur, err := c.Find(ctx, bson.M{})
	if err != nil {
		return ps, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &ps)
	return ps, err
}
