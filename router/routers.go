package router

import (
	"net/http"

	controller "Vegeter/controller"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("POST", "/v1/smartlife/Register", controller.Register, nil)
	register("POST", "/v1/smartlife/AddFriends", controller.AddFriend, nil)
	register("POST", "/v1/smartlife/AddPrice", controller.AddPrice, nil)

	register("POST", "/v1/smartlife/GetCurrentPrice/{uuid}", controller.GetPriceRecord, nil)

}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.Methods(route.Method).
			Path(route.Pattern).
			Handler(route.Handler)
		if route.Middleware != nil {
			r.Use(route.Middleware)
		}
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
