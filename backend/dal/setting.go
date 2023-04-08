package dal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateSetting update existing app configuration
func UpdateSetting(f AppConfig) error {
	c := db.Collection(ColAppConfig)
	ctx := context.Background()

	_, err := c.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{
		"interval": f.Interval,
		"discount": f.Discount,
		"token":    f.Token,
		"chat_id":  f.ChatID,
	}})
	return err
}

// GetAppConfig gets config info
func GetAppConfig() (AppConfig, error) {
	var ac AppConfig
	c := db.Collection(ColAppConfig)
	ctx := context.Background()

	err := c.FindOne(ctx, bson.M{}).Decode(&ac)
	return ac, err
}

// UpdateRunStatus update scraping running status
func UpdateRunStatus(rs int) error {
	c := db.Collection(ColAppConfig)
	ctx := context.Background()

	_, err := c.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{
		"run_scrape": rs,
	}})
	return err
}
