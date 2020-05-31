package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var tmpl *template.Template
var db *sqlx.DB

func main() {
	tmpl = template.Must(template.ParseFiles("templates/search.html"))
	var err error
	db, err = sqlx.Connect("postgres", "user=postgres dbname=polypet_foods sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(db)

	r := mux.NewRouter()
	r.HandleFunc("/", SearchHandler).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

type SearchPageData struct {
	Query        string
	Results      []CatFood
	ResultsCount int
}

type CatFood struct {
	ID          int       `db:"id"`
	Ref         string    `db:"ref"`
	Category    string    `db:"category"`
	Name        string    `db:"name"`
	Ingredients string    `db:"ingredients"`
	Language    string    `db:"language"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

const querySQL = `SELECT category, name, ingredients from cat_foods WHERE ts_idx_col @@ to_tsquery('english', $1)`

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.URL)

	values := r.URL.Query()
	query := values.Get("q")
	var results []CatFood

	err := db.Select(&results, querySQL, query)
	if err != nil {
		log.Println("ERROR", err)
	}
	data := SearchPageData{
		Query:        query,
		Results:      results,
		ResultsCount: len(results),
	}
	w.WriteHeader(http.StatusOK)
	_ = tmpl.Execute(w, data)
}
