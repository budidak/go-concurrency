package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {
	// create a router
	mux := chi.NewRouter()

	// set middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// define application routes
	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/activate-account", app.ActivateAccount)

	// a route and a handler function for sending a mail synchronously
	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "info@mycompany.com",
			FromName:    "info",
			ErrorChan:   make(chan error),
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "Test email",
			Data:    "Hello, World!",
		}

		m.sendMail(msg, make(chan error))
	})

	return mux
}
