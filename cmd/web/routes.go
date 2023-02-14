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
	mux.Get("/activate", app.ActivateAccount)
	mux.Get("/plans", app.chooseSubscription)

	return mux
}
