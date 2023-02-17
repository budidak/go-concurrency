package main

import (
	"context"
	"encoding/gob"
	"go-concurrency/data"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var testApp Config

// TestMain fonksiyonu, testleri çalıştırmaya başladığımızda ilk çalıştırılacak fonksiyondur.
func TestMain(m *testing.M) {
	// setup session
	gob.Register(data.User{})

	tmpPath = "./../../tmp"
	pathToManual = "./../../pdf"

	session := scs.New() // test için yeni bir session oluşturuldu (main.go'daki ile aynı değil)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	// setup config
	testApp = Config{
		Session:       session,
		DB:            nil,
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wait:          &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
		Models:        data.TestNew(nil), // pointer parametre olduğu için nil geçebildik.
		// Artık, testlerimiz için aktif database'e ihtiyaç duymayacak şekilde
		// database fonksiyonlarını tekrar yazdık ve burada Models olarak onları geçtik.
	}

	// create a dummy mailer
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	doneChan := make(chan bool)

	testApp.Mailer = Mail{
		Wait:       testApp.Wait,
		ErrorChan:  errorChan,
		MailerChan: mailerChan,
		DoneChan:   doneChan,
	}

	// fire up goroutines
	go func() {
		for {
			select {
			case <-testApp.Mailer.MailerChan:
				testApp.Wait.Done()
				// mail gönderildiği zaman counter -1 yapıyoruz böylece.
				// Çünkü normal production aşamasında helpers.go içinde sendEmail() ile +1 ekliyorduk.
				// Ama onu yine productiondayken listenForMail() içerisinde mailerChan dinleyerek
				// ne zaman mail gönderilse sendMail() çalıştırıp onun içinde .Done ile -1 yapıyorduk.
				// O listen işlemini burada yapıyoruz.
			case <-testApp.Mailer.ErrorChan:
			case <-testApp.Mailer.DoneChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErrorChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return
			}
		}
	}()

	os.Exit(m.Run()) // setups yapıldıktan sonra test fonksiyonlarını çalıştırır.
}

// getCtx, yazdığımız test fonksiyonlarında request kullandığımızda session bilgisini getirecek.
func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestConfig_SubscribeToPlan(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscribe?id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "user", data.User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Active:    1,
	})

	handler := http.HandlerFunc(testApp.SubscribeToPlan)
	handler.ServeHTTP(rr, req)

	testApp.Wait.Wait() // counter 0 olana kadar bekleriz, çünkü goroutine kullandık.

	// testlerimizi yazalım.
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status code of statusseeother, but got %d", rr.Code)
	}
}
