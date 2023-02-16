package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Bu fonksiyon kısaca session bilgileri doğru ekleniyor mu diye kontrol etmek için yazıldı.
func TestConfig_AddDefaultData(t *testing.T) {
	// AddDefaultData fonksiyonunda r *http.Request kullandığımız için öncelikle bir request oluşturmalı.
	req, _ := http.NewRequest("GET", "/", nil)
	// Daha sonra bu requeste session bilgilerini kabul edebilir hale getirilmeli.
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// test için session bilgilerini ekleyelim
	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	// gerçek fonksiyonumuzu çağıralım (sessiondan bilgi çekiyordu, bakalım doğru çekecek mi)
	td := testApp.AddDefaultData(&TemplateData{}, req)

	if td.Flash != "flash" {
		t.Error("failed to get flash data")
	}

	if td.Warning != "warning" {
		t.Error("failed to get warning data")
	}

	if td.Error != "error" {
		t.Error("failed to get error data")
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	// IsAuthenticated fonksiyonunda r *http.Request kullandığımız için öncelikle bir request oluşturmalı.
	req, _ := http.NewRequest("GET", "/", nil)
	// Daha sonra bu requeste session bilgilerini kabul edebilir hale getirilmeli.
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// gerçek fonksiyonumuzu çağıralım. (sessionda kullanıcı yokken false dönmeli)
	auth := testApp.IsAuthenticated(req)
	if auth {
		t.Error("returns true for authenticated, when it should be false")
	}

	// gerçek fonksiyomuzu çağıralım. (sessiona kullanıcı ekledikten sonra true dönmeli)
	testApp.Session.Put(req.Context(), "userID", 1)
	auth = testApp.IsAuthenticated(req)
	if !auth {
		t.Error("returns false for authenticated, when it should be true")
	}
}

func TestConfig_render(t *testing.T) {
	// zaten cmd/web/ içerisinde $go test -v . çalıştırdığımız için path ona göre yazılmalıdır.
	pathToTemplates = "./templates"

	// render() fonksiyonu http.ResponseWriter ve *http.Request içerdiği için response ve request oluştur.
	rr := httptest.NewRecorder() // response recorder
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// gerçek fonksiyonumuzu çağıralım.
	testApp.render(rr, req, "home.page.gotmpl", &TemplateData{})

	// yukarıdaki fonksiyon çalıştıktan sonra sayfa yüklenemezse 200 kodu dönmez, yani hata vardır.
	if rr.Code != 200 {
		t.Error("failed to render page")
	}
}
