package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dre4success/lenslocked/migrations"
	"github.com/joho/godotenv"

	"github.com/dre4success/lenslocked/controllers"
	"github.com/dre4success/lenslocked/models"
	"github.com/dre4success/lenslocked/templates"
	"github.com/dre4success/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":5050"
	return cfg, nil
}

func main() {
	// set up db connection
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = models.MigrateFs(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// set up services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}

	// set up middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// set up controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	twl := "tailwind.gohtml"
	usersC.Templates.NewTemp = views.Must(views.ParseFS(
		templates.FS, "signup.gohtml", twl))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "signin.gohtml", twl))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS, "forgot-pw.gohtml", twl))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS, "check-your-email.gohtml", twl))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS, "reset-pw.gohtml", twl))

	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "galleries/new.gohtml", twl))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS, "galleries/edit.gohtml", twl))
	galleriesC.Templates.Index = views.Must(views.ParseFS(
		templates.FS, "galleries/index.gohtml", twl))
	galleriesC.Templates.Show = views.Must(views.ParseFS(
		templates.FS, "galleries/show.gohtml", twl))

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

	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
		})
	})

	// start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
