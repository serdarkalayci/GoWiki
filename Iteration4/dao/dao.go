package dao

import (
	"../models"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// WikiDAO Represents Database access values for MongoDBs
type WikiDAO struct {
	Server   string
	Database string
}

var dao *WikiDAO

var db *mgo.Database

// Collection represents MongoDB Collection
const (
	Collection = "wikipages"
)

// Connect method for establishing connection to MongoDB
func (m *WikiDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// SavePage method for saving a new page or updating an existing one
func (m *WikiDAO) SavePage(page *models.Page, newPage bool) error {
	var err error
	if newPage {
		page.ID = bson.NewObjectId()
		err = m.AddNewEntry(page)
	} else {
		err = m.UpdateEntry(page)
	}
	return err
}

// LoadPage method for loading an existing page from database
func (m *WikiDAO) LoadPage(title string) (*models.Page, error) {
	page, err := m.FindByTitle(title)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// AddNewEntry adds new page entry to MongoDB
func (m *WikiDAO) AddNewEntry(page *models.Page) error {
	page.LastUpdate = time.Now()
	err := db.C(Collection).Insert(&page)
	return err
}

// UpdateEntry updates existing page entry on MongoDB
func (m *WikiDAO) UpdateEntry(page *models.Page) error {
	fmt.Printf("Page Id:%v\n", page.ID)
	err := db.C(Collection).Update(bson.M{"_id": page.ID}, bson.M{"title": page.Title, "body": page.Body, "lastUpdate": time.Now()})
	return err
}

// FindByTitle finds a page entry from MongoDB
func (m *WikiDAO) FindByTitle(pageTitle string) (models.Page, error) {
	var page models.Page
	err := db.C(Collection).Find(bson.M{"title": pageTitle}).One(&page)
	return page, err
}

// FindEntries Returns all entries in the database matching the keyword as a slice
func (m *WikiDAO) FindEntries(searchTerm string) (*[]models.Page, error) {
	var pages []models.Page
	err := db.C(Collection).Find(bson.M{"title": bson.M{"$regex": bson.RegEx{searchTerm + "*", ""}, "$options": "i"}}).Sort("-lastUpdate").All(&pages)
	return &pages, err
}

// ListAllEntries Returns all entries in the database as a slice
func (m *WikiDAO) ListAllEntries() (*[]models.Page, error) {
	var pages []models.Page
	err := db.C(Collection).Find(bson.M{}).Sort("-lastUpdate").All(&pages)
	return &pages, err
}
