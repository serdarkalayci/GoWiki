package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/edit.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	dao = &WikiDAO{Server: "127.0.0.1", Database: "gowiki"}
	dao.Connect()
	http.HandleFunc("/", makeHandler(viewHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

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

// AddNewEntry adds new page entry to MongoDB
func (m *WikiDAO) AddNewEntry(page *Page) error {
	err := db.C(Collection).Insert(&page)
	return err
}

// UpdateEntry updates existing page entry on MongoDB
func (m *WikiDAO) UpdateEntry(page *Page) error {
	fmt.Printf("Page Id:%v\n", page.ID)
	err := db.C(Collection).Update(bson.M{"_id": page.ID}, bson.M{"title": page.Title, "body": page.Body})
	return err
}

// FindByTitle finds a page entry from MongoDB
func (m *WikiDAO) FindByTitle(pageTitle string) (Page, error) {
	var page Page
	err := db.C(Collection).Find(bson.M{"title": pageTitle}).One(&page)
	return page, err
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) func(w http.ResponseWriter, p *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/view/deneme"
		}
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	pageID := r.FormValue("id")
	var p *Page
	if pageID != "" {
		id := bson.ObjectIdHex(pageID)
		p = &Page{ID: id, Title: title, Body: []byte(body)}
	} else {
		p = &Page{Title: title, Body: []byte(body)}
	}

	err := p.savePage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) savePage() error {
	var err error
	if !bson.ObjectId(p.ID).Valid() {
		p.ID = bson.NewObjectId()
		err = dao.AddNewEntry(p)
	} else {
		err = dao.UpdateEntry(p)
	}
	return err
}

func loadPage(title string) (*Page, error) {
	// filename := "data/" + title + ".txt"
	// body, err := ioutil.ReadFile(filename)
	page, err := dao.FindByTitle(title)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Page represents a web page to display its title and body
type Page struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Title string        `bson:"title" json:"title"`
	Body  []byte        `bson:"body" json:"body"`
}
