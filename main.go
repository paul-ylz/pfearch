package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
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

type CatFoodCategory struct {
	DisplayName string
	DataName    string
	Attrs       string
}

type SearchPageData struct {
	Query        string
	Results      []CatFood
	ResultsCount int
	Categories   []CatFoodCategory
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

const (
	querySQL = `SELECT category, name, ingredients FROM cat_foods 
				WHERE ts_idx_col @@ to_tsquery('english', $1) 
				  AND category IN (`

	deydrated = "dehydrated and freeze dried"
	dry       = "dry"
	raw       = "frozen raw"
	treat     = "treat"
	wet       = "wet"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.URL)

	values := r.URL.Query()
	var queryArgs []interface{}
	queryArgs = append(queryArgs, values.Get("q"))
	query := values.Get("q")
	categories := getSelectedCategories(values)
	sql := querySQL
	for i, category := range categories {
		queryArgs = append(queryArgs, category)
		sql += fmt.Sprintf("$%d", i+2)
		if i == len(categories)-1 {
			sql += ")"
		} else {
			sql += ","
		}
	}
	log.Println("sql: ", sql)
	var results []CatFood

	err := db.Select(&results, sql, queryArgs...)
	if err != nil {
		log.Println("ERROR", err)
	}

	allCatFoodCategories := []CatFoodCategory{
		{"Dehydrated/freeze-fried", deydrated, inputAttrsFor(values, deydrated)},
		{"Dry", dry, inputAttrsFor(values, dry)},
		{"Raw (frozen)", raw, inputAttrsFor(values, raw)},
		{"Treats", treat, inputAttrsFor(values, treat)},
		{"Wet food", wet, inputAttrsFor(values, wet)},
	}
	data := SearchPageData{
		Query:        query,
		Results:      results,
		ResultsCount: len(results),
		Categories:   allCatFoodCategories,
	}
	w.WriteHeader(http.StatusOK)
	_ = tmpl.Execute(w, data)
}

func getSelectedCategories(values url.Values) []interface{} {
	categories := []string{deydrated, dry, raw, treat, wet}
	var want []interface{}
	for _, category := range categories {
		if values.Get(category) == "1" {
			want = append(want, category)
		}
	}
	return want
}

func inputAttrsFor(values url.Values, param string) string {
	if values.Get(param) == "1" {
		return "checked"
	}
	return ""
}
