package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"
    //"strconv"
	//"text/template"
	forum "Forum/utils"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type M map[string]interface{}

// ANCHOR - VARIABLES

var update *forum.Post
var ID int
var CommentID int
var com string
var txt string = "test"
var User string
var UserID int
var filtre string
var username string
var password string
var U *forum.EntryUser
var E *forum.EntryPost
var C *forum.EntryComment
var L *forum.Like

/*func main() {
	// Create a new ServeMux instance to handle HTTP requests
	mux := http.NewServeMux()

	// Register the handler function for the "/" path
	mux.HandleFunc("/", landingPageHandler)

	// Start the HTTP server on port 8888
	fmt.Println("http://localhost:8888")
	fmt.Println("Starting HTTP server on port 8888...")
	http.ListenAndServe(":8888", mux)
}

//funtion to handle the landing information

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/1landingPage.html")
	//tmpl, err := template.ParseFiles("html/instaPage.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}*/


func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))	
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("pics"))))

	// http.HandleFunc("/", signupHandler)
	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/templates/2loginPage.html", loginPageHandler)
	http.HandleFunc("/templates/2Registration.html", registrationPageHandler)
	http.HandleFunc("/templates/3viewCategoryPage.html", viewCategoryPageHandler)
	http.HandleFunc("/templates/4showArticleCategory.html", showArticleCategoryPageHandler)
	http.HandleFunc("/templates/5showPublication&CommentPage.html", showPublicationPageHandler)
	http.HandleFunc("/templates/6createCategoryPage.html", createCategoryPagePageHandler)
	http.HandleFunc("/templates/7createPostPage.html", createPostPagePageHandler)
	http.HandleFunc("/templates/8aprofilPage.html", aprofilPagePageHandler)
	http.HandleFunc("/templates/9accountActivityPage.html", accountActivityPagePageHandler)
	http.HandleFunc("/templates/10anotherUserProfilePage.html", anotherUserProfilePagePageHandler)
	http.HandleFunc("/templates/aPropos.html", aProposHandler)

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
	renderTemplate(w, "1landingPage.html", nil)
}

func aProposHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "aPropos.html", nil)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "2loginPage.html", nil)
}

func registrationPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "2Registration.html", nil)
}

func viewCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "3viewCategoryPage.html", nil)
}

func showArticleCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "4showArticleCategory.html", nil)
}

func showPublicationPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "5showPublication.html", nil)
}

func createCategoryPagePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "6createCategoryPage.html", nil)
}

func createPostPagePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "7createPostPage.html", nil)
}

func aprofilPagePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "8aprofilPage.html", nil)
}

func accountActivityPagePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "9accountActivityPagePageHandler.html", nil)
}

func anotherUserProfilePagePageHandler(w http.ResponseWriter, r *http.Request) {
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

	templates, _ := template.New("home.html").ParseFiles("templates/8aprofilPage.html")
	// TODO - Gérer le bouton disconnect
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "templates/8aprofilPage.html", U)
	case "POST":
		return
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
		// Affichage basique de la page
		templates.ExecuteTemplate(w, "2Registration.html", U)
	case "POST":
		// Récupération des entrées
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		//TODO - implémenter une fonction de hashage pour la sécurité (BONUS)

		// Synchronisation struct/BDD (import entrées to Struct & export entrées to BDD)
		if username != "" {
			forum.StockageStructUser(U, UserID, username, email, password)
			UserID += 1
			update.AddUser(username, email, password, UserID)
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
	templates, _ := template.New("room.html").ParseFiles("templates/2loginPage.html")

	switch r.Method {
	case "GET":
		// Affichage de la page
		templates.ExecuteTemplate(w, "templates/2loginPage.html", U)
	case "POST":
		// Récupération des entrées
		username = r.FormValue("username/mail")
		password = r.FormValue("loginPassword")
		forum.CheckUserExist(U, E, C, UserID, username, password)

		fmt.Println(U)
		fmt.Println("user = ", username, " password = ", password)
		fmt.Println("isLogged = ", U.IsLogged)
		// MAJ de la page
		templates.ExecuteTemplate(w, "templates/2loginPage.html", U)
	}

}

// TODO - ROOM SYSTEM
/*
func Createroom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/createroom.html")
	}
}

func Joinroom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/joinroom.html")
	}
}

// ANCHOR - ROOM (postes)
func Room(w http.ResponseWriter, r *http.Request) {
	// Initialisation de la struct //
	E = &forum.EntryPost{}
	C = &forum.EntryComment{}
	forum.CheckUserExist(U, E, C, UserID, username, password)
	// Ouverture / création de la bdd //
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(id) from posts").Scan(&ID)
	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID, User, filtre)

	// Initialisation du template //
	templates, _ := template.New("room.html").ParseFiles("templates/room.html")
	fmt.Println("E ", E)
	fmt.Println("U ", U)
	switch r.Method {
	case "GET":
		// affichage de base //
		templates.Execute(w, E)
	case "POST":
		// affichage en cas de création de postes //

		// récupération du texte (uniquement si user il y a)
		txt = r.FormValue("poste")
		if txt != "" {
			filtre = r.FormValue("filtre")
			postID, _ := strconv.Atoi(r.PostFormValue("commentRedirect"))
			println(postID)
			C = &forum.EntryComment{}

			forum.InitialisationComment(C, CommentID)
			C.ID[0] = r.FormValue("commentRedirect")

			E.Text = txt

			// Stockage dans la struct //
			forum.StockageStructPost(E, ID, txt, username, filtre)
			update.LikeNewPost(ID, U)

			ID += 1
			E.User[0] = username
			// Stockage dans la BDD //
			update.AddPost(ID, txt, username, filtre)
			update.LikeNewPost(ID, U)

			// vérification des entrées
			fmt.Println("E.Text : ", E.Text,
				"E.Texts : ", E.Texts,
				"E.ID : ", E.ID,
				"E.User : ", E.Users,
				"E.Filter : ", E.Filters)
		}

		//ANCHOR - FILTRE PAR SUJET
		E.Filter[0] = r.FormValue("searchbar")

		// Mise à jour de la page //
		error := templates.ExecuteTemplate(w, "room.html", E)
		if error != nil {
			panic(error)
		}
	}
}

// ANCHOR - ROOM (commentaires)
func RoomComment(w http.ResponseWriter, r *http.Request) {
	L = &forum.Like{}
	// Ouverture / création de la bdd //
	db, _ := sql.Open("sqlite3", "data/8Forum.db")
	update = forum.NewDB(db)
	U = &forum.EntryUser{}

	update.DB.QueryRow("SELECT MAX(ID) from users").Scan(&UserID)
	forum.InitialisationUsers(U, UserID)
	forum.MajUsers(U, update, UserID, U.Email, U.Password, U.Username)

	// Synchronisation struct/BDD (import from BDD)
	update.DB.QueryRow("SELECT MAX(id) from posts").Scan(&ID)

	forum.InitialisationPost(E, ID)
	forum.MajPost(E, update, ID, User, filtre)

	update.DB.QueryRow("SELECT MAX(CommentID) from comment").Scan(&CommentID)
	forum.InitialisationComment(C, CommentID)
	forum.MajComment(C, update, CommentID, User)

	ID, _ = strconv.Atoi(r.FormValue("commentRedirect"))

	if ID != 0 {
		C = forum.LienCommentPost(C, E, ID)
	}
	update.IsLiked(E, username, ID)
	C.PostIsLiked = E.IsLiked

	C.PostLikes = forum.NbLikesPost(ID, U, update)
	println("likes : ", C.PostLikes)
	templates, _ := template.New("roomComment.html").ParseFiles("templates/roomComment.html")
	switch r.Method {
	case "GET":
		// affichage de base //

		templates.Execute(w, C)
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

		// affichage en cas de création de postes //
		// récupération du texte
		com = r.FormValue("commentaire")
		E.Text = com
		if com != "" {
			forum.StockageStructComment(C, ID, CommentID, com, username)
			CommentID += 1
			update.AddComment(C.ID[0], CommentID, com, username)
			fmt.Println("C.ID = ", C.ID)
			fmt.Println("E.ID = ", E.ID)
			fmt.Println("C.POSTTEXT = ", C.PostText)

		} else {
			if like == "false" {
				E.IsLiked = false
				C.PostIsLiked = E.IsLiked
				update.AddLike(ID, username)
				templates.Execute(w, C)
			} else if like == "true" {
				E.IsLiked = true
				C.PostIsLiked = E.IsLiked
				update.RemoveLike(ID, username)
				templates.Execute(w, C)
			}
		}
		println("like : ", like)

	}
	fmt.Println("C.ID = ", C.ID)
	println("is liked : ", E.IsLiked)

	// Stockage dans la BDD //

}*/
