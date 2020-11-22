package models

import (
	"regexp"
	"unicode"

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
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Category     string             `bson:"category,omitempty"`
	Ingredients  []string           `bson:"ingredients,omitempty"`
	PrepTime     string             `bson:"prepTime,omitempty"`
	Instructions []string           `bson:"instructions,omitempty"`
	OwnerID      string             `bson:"ownerID,omitempty"`
}

// User struct to hold user information
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         string             `bson:"user_id,omitempty"`
	Name           string             `bson:"name,omitempty"`
	Email          string             `bson:"email,omitempty"`
	HashedPassword string             `bson:"hashedPassword,omitempty"`
	Authenticated  bool               `bson:"auth,omitempty"`
	Errors         map[string]string  `bson:"error,omitempty"`
}

// Support struct to pass variable to support template
type Support struct {
	ID      primitive.ObjectID `bson:"id,omitempty"`
	Email   string             `bson:"email,omitempty"`
	Subject string             `bson:"subject,omitempty"`
	Message string             `bson:"message,omitempty"`
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// Validate user emails
func (usr *User) Validate() bool {
	usr.Errors = make(map[string]string)
	match := rxEmail.Match([]byte(usr.Email))
	if match == false {
		usr.Errors["Email"] = "Please enter a valid email address"
	}

	return len(usr.Errors) == 0
}

// Password function that ensures rigorous passwords
func (usr *User) Password(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		usr.Errors["Password"] = "Please make sure your password includes at least [A-Z], [0-9], and symbols."
		return false
	}

	return true
}
