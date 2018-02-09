package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
	"gopkg.in/boj/redistore.v1"
	"github.com/gorilla/sessions"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
const useRedis = true

func App() *buffalo.App {
	if app == nil {
		var sessionStore sessions.Store
		if useRedis {
			store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("secret"))
			if err != nil {
				logrus.Fatalf("Error connecting to redis: %v", err)
			}
			sessionStore = store
		} else {
			sessionStore = sessions.NewCookieStore([]byte("secret"))
		}

		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionName:  "_sessionthing_session",
			SessionStore: sessionStore,
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.GET("/flash", FlashHandler)
		app.ServeFiles("/assets", assetsBox)
	}

	return app
}
