package dal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// collection list
const (
	ColUser      = "users"
	ColSearch    = "searches"
	ColProduct   = "products"
	ColStandard  = "standards"
	ColAppConfig = "app_config"
	ColProxy     = "proxies"
)

// User is for user model
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"`
}

// Search is for search model
type Search struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	URL        string             `json:"url,omitempty" bson:"url,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Postal     string             `json:"postal,omitempty" bson:"postal,omitempty"`
	Status     int                `json:"status"` // 0:deactive, 1:active
	UserID     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	LastSearch *time.Time         `json:"last_search,omitempty" bson:"last_search,omitempty"`
	Items      int                `json:"items,omitempty"`
}

// Product is for product model
type Product struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Asin     string             `json:"asin,omitempty"`
	Title    string             `json:"title,omitempty"`
	URL      string             `json:"url,omitempty" bson:"url,omitempty"`
	Price    float64            `json:"price,omitempty"`
	Prime    float64            `json:"prime,omitempty"`
	SearchID primitive.ObjectID `json:"search_id,omitempty" bson:"search_id,omitempty"`
}

// Standard is for standard product model
type Standard struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Asin  string             `json:"asin,omitempty"`
	Price float64            `json:"price,omitempty"`
	Date  *time.Time         `json:"date,omitempty" bson:"date,omitempty"`
}

// AppConfig is for app config model
type AppConfig struct {
	Interval  uint    `json:"interval,omitempty"`
	Discount  float64 `json:"discount,omitempty"`
	Token     string  `json:"token,omitempty"`
	ChatID    string  `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	RunScrape int     `json:"run_scrape,omitempty" bson:"run_scrape,omitempty"`
}

// Proxy for proxy string model
type Proxy struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Proxy string             `json:"proxy,omitempty"`
}
