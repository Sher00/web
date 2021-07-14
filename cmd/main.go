package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, Fulltext string
}

var posts = []Article{}

type User struct {
	Id                    uint16
	Name, Email, Hashpass string
}

var users = []User{}

func singin(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/singin.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmp.ExecuteTemplate(w, "singin", nil)
}
func singup(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/singup.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmp.ExecuteTemplate(w, "singup", nil)
}
func index(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Fulltext)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	tmp.ExecuteTemplate(w, "index", posts)
}
func save_art(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fulltext := r.FormValue("fulltext")

	if title == "" || anons == "" || fulltext == "" {
		fmt.Fprintf(w, "ne vse dannie")
	} else {
		db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles`(`title`, `anons`, `fulltext`) VALUES('%s', '%s', '%s')", title, anons, fulltext))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/add.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmp.ExecuteTemplate(w, "add", nil)
}

func createuser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	hashpass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if name == "" || email == "" || pass == "" {
		fmt.Fprintf(w, "ne vse dannie")
	} else {
		db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `user`(`name`, `email`, `hashpass`) VALUES('%s', '%s', '%s')", name, email, hashpass))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	pass := r.FormValue("pass")
	//hashpass, err := bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)
	//if err != nil {
	//	fmt.Fprintf(w, err.Error())
	//	}

	if email == "" || pass == "" {
		fmt.Fprintf(w, "ne vse dannie")
	} else {
		db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		search, err := db.Query("SELECT * FROM `user`")
		if err != nil {
			panic(err)
		}
		defer search.Close()

		for search.Next() {
			var us User
			err = search.Scan(&us.Id, &us.Name, &us.Email, &us.Hashpass)
			if err != nil {
				panic(err)
			}
			users = append(users, User{})
			if email == us.Email {
				err = bcrypt.CompareHashAndPassword([]byte(us.Hashpass), []byte(pass))
				if err != nil {
					panic(err)
				} else {
					tmp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

					if err != nil {
						fmt.Fprintf(w, err.Error())
					}
					tmp.ExecuteTemplate(w, "header", users)
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}

			}
		}

	}

}

func handleFunc() {
	http.HandleFunc("/singin/", singin)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/singup/", singup)
	http.HandleFunc("/createuser/", createuser)
	http.HandleFunc("/add/", add)
	http.HandleFunc("/save_art/", save_art)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func main() {

	handleFunc()
}
