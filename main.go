package main

import (
	"fmt"
	"github.com/dre4success/lenslocked/migrations"
	"net/http"

	"github.com/dre4success/lenslocked/controllers"
	"github.com/dre4success/lenslocked/models"
	"github.com/dre4success/lenslocked/templates"
	"github.com/dre4success/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	// set up db connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = models.MigrateFs(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// set up services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)

	// set up controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	twl := "tailwind.gohtml"
	usersC.Templates.NewTemp = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", twl))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", twl))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS, "forgot-pw.gohtml", twl))

	// set up router and routes
	r := chi.NewRouter()
	// middleware to be used everywhere
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS,
		"home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", twl))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", twl))))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Post("/signout", usersC.ProcessSignOut)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	// start the server
	fmt.Println("Starting the server on :5050...")
	err = http.ListenAndServe(":5050", r)
	if err != nil {
		panic(err)
	}
}
