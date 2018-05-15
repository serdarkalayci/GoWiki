package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Page represents a web page to display its title and body
type Page struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Title      string        `bson:"title" json:"title"`
	Body       string        `bson:"body" json:"body"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
	Comments   []Comment     `bson:"comments" json:"comments"`
}
