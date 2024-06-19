package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	forum "Forum/utils"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type M map[string]interface{}

// ANCHOR - VARIABLES

var update *forum.Post
var ID int
var CommentID int
var User string
var UserID int
var filtre string
var username string
var password string
var U *forum.EntryUser
var E *forum.EntryPost
var C *forum.EntryComment
var L *forum.Like

var db *sql.DB

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("pics"))))

	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/registration", Register)
	http.HandleFunc("/categories", viewCategoryPageHandler)
	http.HandleFunc("/category", showArticleCategoryPageHandler)
	http.HandleFunc("/publications", showPublicationPageHandler)
	http.HandleFunc("/create-category", createCategoryPagePageHandler)
	http.HandleFunc("/create-post", createPostPageHandler)
	http.HandleFunc("/profil-self", aprofilPagePageHandler)
	http.HandleFunc("/activity", accountActivityPageHandler)
	http.HandleFunc("/profil", anotherUserProfilePageHandler)
	http.HandleFunc("/about", aProposHandler)
	http.HandleFunc("/disconnect", disconnect)

	fmt.Println("http://localhost:8888/")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	Home(w, r)
}

func aProposHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "aPropos.html", nil)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	Login(w, r)
	// renderTemplate(w, "2loginPage.html", nil)
}

func disconnect(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "userId",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "../", http.StatusSeeOther)
}

func viewCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "3viewCategoryPage.html", nil)
}

func showArticleCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "4showArticleCategory.html", nil)
}

func showPublicationPageHandler(w http.ResponseWriter, r *http.Request) {
	E := &forum.EntryPost{}

	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID, E.User[ID], E.Filter[ID])
	templates, _ := template.New("5showPublication&CommentPage.html").ParseFiles("templates/5showPublication&CommentPage.html")
	templates.ExecuteTemplate(w, "5showPublication&CommentPage.html", E)
}

func createCategoryPagePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "6createCategoryPage.html", nil)
}

func createPostPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Redirect(w, r, "../login", http.StatusSeeOther)
		return
	}
	templates, _ := template.New("7createPostPage.html").ParseFiles("templates/7createPostPage.html")
	switch r.Method {
	case "GET":
		println("On new post")
		// Affichage basique de la page
		templates.ExecuteTemplate(w, "7createPostPage.html", U)
	case "POST":
		// Récupération des entrées
		text := r.FormValue("post-content")
		user := cookie.Value
		filter := r.FormValue("post-category")

		println("content : " + text + "\nauthor : " + user + "\ncategory : " + filter)

		//TODO - implémenter une fonction de hashage pour la sécurité (BONUS)

		// Synchronisation struct/BDD (import entrées to Struct & export entrées to BDD)
		if username != "" {
			forum.StockageStructPost(E, ID, text, user, filter)
			ID += 1
			update.AddPost(ID, text, user, filter)
			println("E : ")
			fmt.Println(E)
		}
		templates.ExecuteTemplate(w, "7createPostPage.html", U)
	}

}

func aprofilPagePageHandler(w http.ResponseWriter, r *http.Request) {
	U = &forum.EntryUser{}
	E = &forum.EntryPost{}
	C = &forum.EntryComment{}

	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Redirect(w, r, "../login", http.StatusSeeOther)
		return
	}
	for _, cookie := range r.Cookies() {
		fmt.Println(w, cookie.Name)
	}
	U.ID, _ = strconv.Atoi(cookie.Value)
	println("ID : ")
	println(U.ID)

	update.DB.QueryRow("SELECT email, username, password  from users WHERE ID = "+cookie.Value).Scan(&U.Email, &U.Username, &U.Password)
	println("Username : " + U.Username)
	templates, _ := template.New("8aprofilPage.html").ParseFiles("templates/8aprofilPage.html")
	templates.ExecuteTemplate(w, "8aprofilPage.html", U)
}

func accountActivityPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "9accountActivityPage.html", nil)
}

func anotherUserProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "10anotherUserProfilePage.html", nil)
}

// ANCHOR - PAGE QUI SAIT ACCEUILLIR
func Home(w http.ResponseWriter, r *http.Request) {
	U = &forum.EntryUser{}
	E = &forum.EntryPost{}
	C = &forum.EntryComment{}

	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	update.DB.QueryRow("SELECT MAX(ID) from users").Scan(&UserID)
	forum.InitialisationUsers(U, UserID)
	forum.MajUsers(U, update, UserID, U.Email, password, username)
	forum.CheckUserExist(U, E, C, UserID, username, password)

	update.DB.QueryRow("SELECT MAX(ID) from posts").Scan(&ID)
	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID, User, filtre)
	update.LikeUser(U, E)

	templates, _ := template.New("1landingPage.html").ParseFiles("templates/1landingPage.html")
	// TODO - Gérer le bouton disconnect
	switch r.Method {
	case "GET":
		// for _, cookie := range r.Cookies() {
		// 	println(w, cookie.Name)
		// }
		templates.ExecuteTemplate(w, "1landingPage.html", U)
	}
}

// ANCHOR - PAGE INSCRIPTION

func Register(w http.ResponseWriter, r *http.Request) {
	// Ouverture / création de la bdd //
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	templates, _ := template.New("2Registration.html").ParseFiles("templates/2Registration.html")

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(ID) from users").Scan(&UserID)
	forum.InitialisationUsers(U, UserID)
	forum.MajUsers(U, update, UserID, U.Email, U.Password, U.Username)

	update.DB.QueryRow("SELECT MAX(ID) from posts").Scan(&ID)
	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID, User, filtre)

	forum.CheckUserExist(U, E, C, UserID, username, password)

	//TODO - faire que le bouton inscription renvoi les data en plus d'envoyer sur la page d'accueil (il renvoi sur une page blanche)
	switch r.Method {
	case "GET":
		println("On register")
		// Affichage basique de la page
		templates.ExecuteTemplate(w, "2Registration.html", U)
	case "POST":
		// Récupération des entrées
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		println("name : " + username + "\nemail : " + email + "\npassword : " + password)

		//TODO - implémenter une fonction de hashage pour la sécurité (BONUS)

		// Synchronisation struct/BDD (import entrées to Struct & export entrées to BDD)
		if username != "" {
			forum.StockageStructUser(U, UserID, username, email, password)
			UserID += 1
			update.AddUser(username, email, password, UserID)
			println("U : ")
			fmt.Println(U)
			update.LikeNewUser(username, E)
			U.IsSignedIn = true
		}
		templates.ExecuteTemplate(w, "2Registration.html", U)
	}

}

// ANCHOR - LOGIN
func Login(w http.ResponseWriter, r *http.Request) {
	// Initialisation de la struct //
	U = &forum.EntryUser{IsLogged: false}

	// Ouverture / création de la bdd //
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(ID) from users").Scan(&UserID)
	forum.InitialisationUsers(U, UserID)
	forum.MajUsers(U, update, UserID, U.Email, U.Password, U.Username)

	// Initialisation du template //
	templates, _ := template.New("2loginPage.html").ParseFiles("templates/2loginPage.html")
	println("login operations")
	switch r.Method {
	case "GET":
		// Affichage de la page
		templates.ExecuteTemplate(w, "2loginPage.html", U)
	case "POST":
		// Récupération des entrées
		username = r.FormValue("username/mail")
		password = r.FormValue("loginPassword")
		forum.CheckUserExist(U, E, C, UserID, username, password)

		fmt.Println(U)
		fmt.Println("user = ", username, " password = ", password)
		fmt.Println("isLogged = ", U.IsLogged)

		if U.IsLogged {
			cookie := http.Cookie{
				Name:    "userId",
				Value:   strconv.Itoa(U.ID + 1),
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "../profil-self", http.StatusSeeOther)
		} else {
			templates.ExecuteTemplate(w, "2loginPage.html", U)
		}

		// MAJ de la page

	}

}
