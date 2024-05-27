package main

import (
	"fmt"
	"net/http"

	"github.com/dre4success/lenslocked/controllers"
	"github.com/dre4success/lenslocked/models"
	"github.com/dre4success/lenslocked/templates"
	"github.com/dre4success/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	r := chi.NewRouter()

	twl := "tailwind.gohtml"
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS,
		"home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", twl))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", twl))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.NewTemp = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", twl))

	r.Get("/signup", usersC.New)

	r.Post("/signup", usersC.Create)

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", twl))

	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	r.Get("/users/me", usersC.CurrentUser)
	r.Post("/signout", usersC.ProcessSignOut)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)
	fmt.Println("Starting the server on :5050...")
	err = http.ListenAndServe(":5050", csrfMw(r))
	if err != nil {
		panic(err)
	}
}
