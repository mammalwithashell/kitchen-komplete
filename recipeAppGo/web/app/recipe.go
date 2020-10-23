package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe struct for mongodb
type Recipe struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Class string             `bson:"class,omitempty"`
}
