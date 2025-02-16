package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

type Article struct {
	Title string
	Body  string
}

var prestr string

func main() {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := "CREATE TABLE IF NOT EXISTS article(title STRING,body STRING)"
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := "SELECT * FROM article"
	row := DbConnection.QueryRow(cmd)
	var a Article
	err := row.Scan(&a.Title, &a.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No row")
		} else {
			log.Println(err)
		}
	}
	renderTemplate(w, "view", &a)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := "SELECT * FROM article"
	row := DbConnection.QueryRow(cmd)
	var a Article
	err := row.Scan(&a.Title, &a.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No row")
		} else {
			log.Println(err)
		}
	}
	prestr = a.Body
	renderTemplate(w, "edit", &a)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	DbConnection, _ := sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	if body != prestr {
		cmd := "INSERT INTO article (title,body) VALUES (?,?)"
		_, err := DbConnection.Exec(cmd, "Go言語日記", body)
		if err != nil {
			log.Fatalln(err)
		}
		cmd = "DELETE FROM article WHERE body = ?"
		_, err = DbConnection.Exec(cmd, prestr)
		if err != nil {
			log.Fatalln(err)
		}
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, a *Article) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, a)
}
