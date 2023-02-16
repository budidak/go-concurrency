package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// NOTE: run all tests with "$go test -v ." while inside the cmd/web/ directory.

// projede yazdığımız routes tekrar yazdık.
var routes = []string{
	"/",
	"/login",
	"/logout",
	"/register",
	"/activate",
	"/members/plans",
	"/members/subscribe",
}

// Test_Routest_Exist function tests if the each route in the routes are actually exists in the project
func Test_Routes_Exist(t *testing.T) {
	testRoutes := testApp.routes()

	chiRoutes := testRoutes.(chi.Router) // http.Handler -> chi.Router typecast yaptık

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	// herhangi bir router tree üzerinde gezinebilmek için chi.Walk() kullanırız.
	_ = chi.Walk(
		routes,
		// ikinci parametremiz bir fonksiyon, routes üzerinde gezinirken bu fonk. argümanlarını alır.
		func(method string, foundRoute string, handler http.Handler,
			middlewares ...func(http.Handler) http.Handler) error {
			if route == foundRoute {
				found = true
			}
			return nil
		})

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
