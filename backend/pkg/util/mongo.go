package util

import "go.mongodb.org/mongo-driver/bson/primitive"

// ConvStrToObjID converts string to mongo object id
func ConvStrToObjID(s string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}
