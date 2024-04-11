package main

import (
	"fmt"
	"net/http"

	"github.com/dre4success/lenslocked/controllers"
	"github.com/dre4success/lenslocked/templates"
	"github.com/dre4success/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	twl := "tailwind.gohtml"
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS,
		"home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", twl))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", twl))))

	var usersC controllers.Users
	usersC.Templates.NewTemp = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", twl))

	r.Get("/signup", usersC.New)

	r.Post("/signup", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :5050...")
	err := http.ListenAndServe(":5050", r)
	if err != nil {
		panic(err)
	}
}
