package supermarket

import "go.mongodb.org/mongo-driver/bson/primitive"

type Supermarket struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Address   string             `json:"address" bson:"address"`
	Latitude  float64            `json:"latitude" bson:"latitude"`
	Longitude float64            `json:"longitude" bson:"longitude"`
}
