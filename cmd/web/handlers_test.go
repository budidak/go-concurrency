package main

import (
	"go-concurrency/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Birden fazla page olduğu için benzer işleri yapacağımızdan dolayı bir table oluşturduk,
// bu table üzerinde döngüyle gezinerek her sayfa için aynı olan testlerimizi çalıştıracağız.
var pageTests = []struct {
	name               string           // sayfanın adı
	url                string           // sayfanın adresi
	expectedStatusCode int              // sayfada beklenilen http statuscode
	handler            http.HandlerFunc // sayfa ile ilgili handler function
	sessionData        map[string]any   // sayfada session bilgisi kullanılmışsa
	expectedHTML       string           // sayfada görmeyi umduğumuz bir html varsa
}{
	{ // örneğin home page yüklendiğinde aşağıdaki gibi şeyleri umuyoruz.
		name:               "home",
		url:                "/",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.HomePage,
	},
	{ // login page yüklendiğinde aşağıdaki gibi şeyleri umuyoruz.
		name:               "login",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{ // logout page yönlendirme yapıyor ve session kullanıyor, o yüzden aşağıdaki gibi yazdık.
		name:               "logout",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.Logout,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
	{ // register page içerisinde aşağıdaki gibi şeyler görmeyi umuyoruz.
		name:               "register",
		url:                "/register",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.RegisterPage,
		expectedHTML:       `<h1 class="mt-5">Register</h1>`,
	},
}

func Test_Pages(t *testing.T) {
	pathToTemplates = "./templates" // kodu çalıştırırken zaten cmd/web/ klasörünün içerisindeyiz.

	// öncelikle tabloya eklediğimiz tüm sayfalar üzerinde döngü kurduk.
	for _, e := range pageTests {
		rr := httptest.NewRecorder()                 // response
		req, _ := http.NewRequest("GET", e.url, nil) // request
		ctx := getCtx(req)
		req = req.WithContext(ctx) // request that can accept session data

		// eğer session data'mız varsa bunu context içerisinde alalım.
		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		// handleri kullanarak sayfayı serve et.
		e.handler.ServeHTTP(rr, req)

		// teslerimizi yazıyoruz.
		// eğer status code beklenilenden farklıysa hata vardır:
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s failed: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		// eğer sayfa render edildiğinde görmeyi umduğumuz html'i göremezsek hata vardır:
		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("%s failed: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}
