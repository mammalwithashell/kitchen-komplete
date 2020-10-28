package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Pantry object for mongodb
type Pantry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id,omitempty"`
	Ingredients []string           `bson:"ingredients,omitempty"`
}

// Recipe struct for mongodb
type Recipe struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Class string             `bson:"class,omitempty"`
}
