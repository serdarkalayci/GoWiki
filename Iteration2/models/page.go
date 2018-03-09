package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Page represents a web page to display its title and body
type Page struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Title string        `bson:"title" json:"title"`
	Body  []byte        `bson:"body" json:"body"`
}
