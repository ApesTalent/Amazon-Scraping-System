package dal

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSearch creates new search
func CreateSearch(s Search) (Search, error) {
	c := db.Collection(ColSearch)
	ctx := context.Background()

	s.URL = strings.TrimSpace(s.URL)
	rst, err := c.InsertOne(ctx, s)
	s.ID = rst.InsertedID.(primitive.ObjectID)

	log.Println(s)
	return s, err
}

// UpdateSearch updates existing search
func UpdateSearch(s Search) error {
	c := db.Collection(ColSearch)
	ctx := context.Background()

	s.URL = strings.TrimSpace(s.URL)
	_, err := c.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{
		"url":         s.URL,
		"last_search": s.LastSearch,
		"status":      s.Status,
	}})
	return err
}

// DeleteSearch remove one search data
func DeleteSearch(id primitive.ObjectID) error {
	c := db.Collection(ColSearch)
	ctx := context.Background()

	_, err := c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListSearch get all searches
func ListSearch() ([]Search, error) {
	ss := []Search{}
	c := db.Collection(ColSearch)
	ctx := context.Background()

	cur, err := c.Find(ctx, bson.D{})
	if err != nil {
		return ss, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &ss)
	if err != nil {
		return ss, err
	}

	for i := 0; i < len(ss); i++ {
		ps, _ := ListProduct(ss[i].ID)
		ss[i].Items = len(ps)
	}
	return ss, err
}

// GetSearchByID get a search by id
func GetSearchByID(id primitive.ObjectID) (*Search, error) {
	s := &Search{}
	c := db.Collection(ColSearch)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	return s, err
}
