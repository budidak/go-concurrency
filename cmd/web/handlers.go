package main

import (
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gotmpl", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gotmpl", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	// we need to renew token each time the user login/logout
	_ = app.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err) // for development phase
	}

	// get email and password from the post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// get user from database by email
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		// eğer parola yanlış girilirse serverdan otomatik bir mail yollat.
		msg := Message{
			To:      email,
			Subject: "Failed to login attempt!",
			Data:    "Invalid login attempt!",
		}
		app.sendEmail(msg) // bu fonksiyon channel'a msg parametresini iletecek.
		// channela bir data geldiğinde arka planda goroutine çalıştığı için
		// server da otomatik olarak sistemden mail yollayacak.
		app.Session.Put(r.Context(), "error", "Invalid credentials...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// successfully logged in;
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Successful login!")

	// redirect the user to the homepage after logged in
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gotmpl", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	// create a user

	// send an activation email (processing time & delay açısından expensive olabilir o yüzden concurrency ile arkaplanda halledeceğiz.)
	// email ile gönderdiğimiz link tıklandığında aşağıdaki fonksiyonu çalıştıracak.

	// subscribe the user to an account
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session and renew the token
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}
