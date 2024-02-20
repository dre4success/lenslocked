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

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS,
		"home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :7070...")
	http.ListenAndServe(":7070", r)
}
