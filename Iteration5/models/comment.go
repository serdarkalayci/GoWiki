package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Comment represents a comment entered for a page or another comment
type Comment struct {
	CommentID   bson.ObjectId `bson:"commentId" json:"commentId"`
	Body        string        `bson:"body" json:"body"`
	CommentDate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
	Comments    []Comment     `bson:"comments" json:"comments"`
}
