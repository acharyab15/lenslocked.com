package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>OH NO!</h1> The URL you were trying to visit wasn't found!!")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	usersC := controllers.NewUsers()
	staticC := controllers.NewStatic()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	var h http.Handler = http.HandlerFunc(notFound)
	r.NotFoundHandler = h
	http.ListenAndServe(":3000", r)
}
