package main

import (
	"fmt"
	"html/template"
	"net/http"

	"totp/types"
	"totp/users"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/time", mainHandler)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/signup", users.Signup)
	http.HandleFunc("/qrcode", users.QRCode)
	http.HandleFunc("/validate", users.Validate)

	img := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", img))

	portStr := fmt.Sprintf("localhost:%d", 62222)
	fmt.Println("Listening on:", portStr)
	http.ListenAndServe(portStr, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var (
		header  types.HeaderRecord
		lir     types.LoginRec
		theList []types.LoginRec
	)

	if !users.LoggedIn(w, r, &lir) {
		users.DisplayWelcome(w, r)
		return
	}

	c, _ := r.Cookie("session")

	theUser := users.GetLogin(c.Value)
	header.Title = "Edinburgh Go users"
	data := struct {
		Header  types.HeaderRecord
		Name    string
		TheList []types.LoginRec
	}{
		header,
		theUser.Mail,
		theList,
	}

	t, err := template.ParseFiles("templates/index.html", types.ViewHeader, types.ViewMenuUs)
	if err != nil {
		fmt.Println("main page err", err)
	}
	t.Execute(w, data)
}
