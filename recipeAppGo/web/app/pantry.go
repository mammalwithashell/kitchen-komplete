package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Pantry object for mongodb
type Pantry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id,omitempty"`
	Ingredients []string           `bson:"ingredients,omitempty"`
}
