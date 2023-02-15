package main

import (
	"fmt"
	"go-concurrency/data"
	"html/template"
	"net/http"
	"strconv"
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
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// TODO - Validate data (kullanıcının girdiği veriler doğru mu, kullanıcı daha önceden kayıtlı mı vs.)

	// create a user
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	// normalde env.var.'dan almamız lazım burada basitçe hardcoded şekilde yazdık.
	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(url) // prevents url tampering
	app.InfoLog.Println(signedURL)            // just printed it to see

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL), // any olduğu için herhangi bir tipte veri koyabildik.
	}

	app.sendEmail(msg)
	app.Session.Put(r.Context(), "flash", "Confirmation email sent. Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session and renew the token
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url
	url := r.RequestURI                               // get url from request (clicked activation link)
	testURL := fmt.Sprintf("http://localhost%s", url) // env.var'dan almalıyız basitçe hardcoded yazdık.
	okay := VerifyToken(testURL)
	if !okay {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate account, url email parametresine göre db'den ilgili user'i çek.
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1 // db'den user başarıyla çekildi burası çalıştı.
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// update kısmında da sorun olmazsa buradayız.
	app.Session.Put(r.Context(), "flash", "Account activated! You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// send an email with attachments

}

func (app *Config) chooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err) // terminale yazmak yerine, error sayfası oluşturup yönlendirebilirsin.
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gotmpl", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is chosen (bunu URL üzerinden alabiliriz, plans.page.gotmpl yazdık)
	id := r.URL.Query().Get("id")
	planID, _ := strconv.Atoi(id)

	// get the plan from the database (URL'den alınan id bilgisi ile veritabanından planı getir)
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find plan.")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	// get the user from the session (has user logged in or not?)
	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Log in first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice (concurrent task)

	// send an email with the invoice attached (concurrent task)

	// generate a manual (concurrent task)

	// send an email with the manual attached (concurrent task)

	// subscribe the user to an account

	// redirect

}
