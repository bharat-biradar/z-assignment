package application

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) Router() http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router := httprouter.New()
	middlewareChain := alice.New(app.recoverPanic, app.logRequest, app.setHeaders)
	authorized := alice.New(app.isAuthorized)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/api/items", authorized.ThenFunc(app.getAllItems))
	router.Handler(http.MethodGet, "/api/item/:id", authorized.ThenFunc(app.getItem))
	router.Handler(http.MethodPost, "/api/items", authorized.ThenFunc(app.addItem))
	router.Handler(http.MethodDelete, "/api/item/:id", authorized.ThenFunc(app.deleteItem))
	router.Handler(http.MethodPatch, "/api/item/:id", authorized.ThenFunc(app.updateItem))

	router.HandlerFunc(http.MethodPost, "/api/user/signup", app.createUser)
	router.HandlerFunc(http.MethodPost, "/api/user/login", app.loginUser)
	router.Handler(http.MethodDelete, "/api/user/delete", authorized.ThenFunc(app.deletUser))

	router.HandlerFunc(http.MethodGet,"/",app.homePage)
	return middlewareChain.Then(router)
}
