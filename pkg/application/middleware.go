package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"task1/items_manager/pkg/models"
)

type contextKey string

const contextUserKey contextKey = "username"

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				if errors.Is(err.(error), models.ErrInvalidApiKey) {
					app.clientError(w, http.StatusUnauthorized)
					return
				}

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		// w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Content-Security-Policy",
		// 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		// w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *Application) isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api_key := r.Header.Get("api_key")
		// app.infoLog.Printf()
		// fmt.Println("api key",api_key)
		val, ok := app.sessionManager.isTokenValid(api_key)
		if !ok {
			// fmt.Println("api key invalid")
			app.sessionManager.removeSession(val)

			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.infoLog.Printf("%s-%s", val.Username, val.Token)
		// fmt.Println("vaalid key User ", val.Username)
		ctx := context.WithValue(r.Context(), contextUserKey, val.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
