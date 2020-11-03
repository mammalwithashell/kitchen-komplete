package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pantry object for mongodb
type Pantry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id,omitempty"`
	Ingredients map[string]int     `bson:"ingredients,omitempty"`
}

// Recipe struct for mongodb
type Recipe struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Class       string             `bson:"class,omitempty"`
	Ingredients []string           `bson:"ingredients,omitempty"`
}

// Support struct to pass variable to support template
type Support struct {
	ID      primitive.ObjectID `bson:"id,omitempty"`
	Email   string             `bson:"email,omitempty"`
	Subject string             `bson:"subject,omitempty"`
	Message string             `bson:"message,omitempty"`
}

// User struct to hold user information
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         string             `bson:"user_id,omitempty"`
	Name           string             `bson:"name,omitempty"`
	Email          string             `bson:"email,omitempty"`
	HashedPassword string             `bson:"hashedPassword,omitempty"`
	Authenticated  bool               `bson:"auth,omitempty"`
}
