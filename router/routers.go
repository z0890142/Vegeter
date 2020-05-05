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
	register("POST", "/v1/Register", controller.Register, nil)
	register("POST", "/v1/AddFriends", controller.AddFriend, nil)
	register("POST", "/v1/AddPrice", controller.AddPrice, nil)

	register("PUT", "/v1/ConfirmFriend", controller.ConfirmFriend, nil)

	register("GET", "/v1/GetFriends/{uuid}", controller.GetFriend, nil)
	register("GET", "/v1/AddFriends/{uuid}", controller.GetAddFriend, nil)
	register("GET", "/v1/GetCurrentPrice/{uuid}", controller.GetPriceRecord, nil)
	register("GET", "/v1/GetAllPrice", controller.GetALLRecord, nil)

	register("GET", "/v1/GetUserInfo/{uuid}", controller.GetUserInfo, nil)

	register("GET", "/ws/{room}", controller.StartWebSocket, nil)

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
