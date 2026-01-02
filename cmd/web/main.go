package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yayifications/bookings/pkg/config"
	"github.com/yayifications/bookings/pkg/handlers"
	"github.com/yayifications/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const PORTNUMBER = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the entry point of the application
func main() {

	var app config.AppConfig

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // in production you want it to be true
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    PORTNUMBER,
		Handler: routes(&app),
	}

	fmt.Printf("Starting application on port %s\n", PORTNUMBER)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
