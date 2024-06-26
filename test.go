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
var username string
var password string
var U *forum.EntryUser
var E *forum.EntryPost
var C *forum.EntryComment
var L *forum.Like

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("pics"))))

	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/registration", Register)
	http.HandleFunc("/categories", viewCategoryPageHandler)
	http.HandleFunc("/category", showArticleCategoryPageHandler)
	http.HandleFunc("/publications", showPublicationPageHandler)
	http.HandleFunc("/comment", commentPageHandler)
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
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)
	var categories []string

	rows, err := update.DB.Query("SELECT Name FROM categories;")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		rows.Scan(&category)
		categories = append(categories, category)
	}

	fmt.Println("categories : ", categories)

	templates, _ := template.New("3viewCategoryPage.html").ParseFiles("templates/3viewCategoryPage.html")
	templates.ExecuteTemplate(w, "3viewCategoryPage.html", categories)
}

func showArticleCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	E = &forum.EntryPost{}
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	update.DB.QueryRow("SELECT MAX(ID) from posts").Scan(&ID)
	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID)
	E.Filter = r.URL.Query().Get("category")

	fmt.Println("Struct : ", E)

	fmt.Println("Filter : ", E.Filter)

	templates, _ := template.New("4showArticleCategory.html").ParseFiles("templates/4showArticleCategory.html")

	templates.ExecuteTemplate(w, "4showArticleCategory.html", E)
}

func showPublicationPageHandler(w http.ResponseWriter, r *http.Request) {
	E = &forum.EntryPost{}
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(ID) from posts").Scan(&ID)
	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID)
	fmt.Println("Text : "+E.Text+"\nUser : ", E.User)
	fmt.Println("Struct : ", E)
	templates, _ := template.New("5showPublication&CommentPage.html").ParseFiles("templates/5showPublication&CommentPage.html")
	fmt.Println("E.Filter : ", E.Filter)

	// Affichage basique de la page
	templates.ExecuteTemplate(w, "5showPublication&CommentPage.html", E)

}

func commentPageHandler(w http.ResponseWriter, r *http.Request) {
	L = &forum.Like{}
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)
	U = &forum.EntryUser{}

	update.DB.QueryRow("SELECT MAX(ID) from users").Scan(&UserID)
	forum.InitialisationUsers(U, UserID)
	forum.MajUsers(U, update, UserID, U.Email, U.Password, U.Username)

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(id) from posts").Scan(&ID)

	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID)

	update.DB.QueryRow("SELECT MAX(CommentID) from comment").Scan(&CommentID)
	forum.InitialisationComment(C, CommentID)
	forum.MajComment(C, update, CommentID, User)

	postId, _ := strconv.Atoi(r.URL.Query().Get("postID"))
	C = forum.LienCommentPost(C, E, postId)

	userId, err := r.Cookie("userId")
	if err != nil {
		C.IsLogged = false
	} else {
		C.IsLogged = true
		update.DB.QueryRow("SELECT Username from users WHERE ID = ?", userId.Value).Scan(&username)
		fmt.Println("SELECT Username from users WHERE ID = ", userId.Value)
		fmt.Println("userId : ", userId.Value)
		fmt.Println("USERNAME : ", username)

	}
	id, _ := strconv.Atoi(userId.Value)
	update.IsLiked(E, username, id)
	C.PostIsLiked = E.IsLiked

	C.PostLikes = forum.NbLikesPost(id, U, update)

	fmt.Println("C : ", C)
	fmt.Println("E : ", E)
	fmt.Println("U : ", U)
	templates, _ := template.New("11commentPage.html").ParseFiles("templates/11commentPage.html")

	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "11commentPage.html", C)
	case "POST":

		update.DB.QueryRow("SELECT ID FROM from posts WHERE text = ?", C.PostText).Scan(&ID)
		forum.GetLikesFromPost(L, update, ID)
		for i := 0; i < len(L.User); i++ {
			if L.User[i] == C.PostUser && L.Like[i] == "true" {
				E.IsLiked = true
				C.PostIsLiked = true
			} else {
				E.IsLiked = false
				C.PostIsLiked = false
			}
		}
		C.PostIsLiked = E.IsLiked
		like := r.FormValue("like")

		com := r.FormValue("commentaire")
		if com != "" {
			forum.StockageStructComment(C, postId, CommentID, com, username)
			CommentID += 1
			update.AddComment(C.ID[0], CommentID, com, username)
			fmt.Println("C.ID = ", C.ID)
			fmt.Println("E.ID = ", E.ID)
			fmt.Println("C.POSTTEXT = ", C.PostText)

		} else {
			if like == "false" {
				E.IsLiked = false
				C.PostIsLiked = E.IsLiked
				update.AddLike(postId, username)
				fmt.Println("like added at username " + username)
				templates.Execute(w, C)
			} else if like == "true" {
				E.IsLiked = true
				C.PostIsLiked = E.IsLiked
				update.RemoveLike(postId, username)
				fmt.Println("like removed")
				templates.Execute(w, C)
			}
		}
		println("like : ", like)
	}

}

func createCategoryPagePageHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	var categoryId int

	update.DB.QueryRow("SELECT MAX(ID) from categories").Scan(&categoryId)
	switch r.Method {
	case "POST":
		name := r.FormValue("category-name")
		description := r.FormValue("category-description")
		fmt.Println("name : ", name, " description : ", description)

		stmt, err := update.DB.Prepare(`INSERT INTO categories (ID, Name, Description) VALUES (?, ?, ?)`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stmt.Exec(strconv.Itoa(categoryId+1), name, description)

		http.Redirect(w, r, "../", http.StatusAccepted)
	}
	fmt.Println("Create category")
	renderTemplate(w, "6createCategoryPage.html", nil)

}

func createPostPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Redirect(w, r, "../login", http.StatusSeeOther)
		return
	}

	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	var categories []string
	rows, err := update.DB.Query("SELECT Name from categories")
	for rows.Next() {
		var category string
		rows.Scan(&category)
		categories = append(categories, category)
	}
	U.Categories = categories

	fmt.Println("U.Categories : ", U.Categories)
	templates, _ := template.New("7createPostPage.html").ParseFiles("templates/7createPostPage.html")
	switch r.Method {
	case "GET":
		println("On new post")
		// Affichage basique de la page
		templates.ExecuteTemplate(w, "7createPostPage.html", U)
	case "POST":
		var user string
		// Récupération des entrées
		text := r.FormValue("post-content")
		userid := cookie.Value
		filter := r.FormValue("post-category")

		err := update.DB.QueryRow("SELECT username from users WHERE ID = ?", userid).Scan(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		println("content : " + text + "\nauthor : " + user + "\ncategory : " + filter)

		//TODO - implémenter une fonction de hashage pour la sécurité (BONUS)

		// Synchronisation struct/BDD (import entrées to Struct & export entrées to BDD)
		if user != "" {
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
	forum.MajPost(E, update, ID)
	update.LikeUser(U, E)

	_, err := r.Cookie("userId")
	if err != nil {
		U.IsLogged = false
	} else {
		U.IsLogged = true
	}

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
	forum.MajPost(E, update, ID)

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
