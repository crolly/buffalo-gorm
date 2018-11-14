package actions

import (
	"strconv"

	"github.com/crolly/buffalo-gorm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/pop"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"

	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_gorm_test_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		// app.Use(popmw.Transaction(models.DB))
		app.Use(GormTransaction(models.GormDB))

		// Setup and use translations:
		app.Use(translations())
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

var GormTransaction = func(db *gorm.DB) buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

			ef := func() error {
				if err := h(c); err != nil {
					return err
				}
				if res, ok := c.Response().(*buffalo.Response); ok {
					if res.Status < 200 || res.Status >= 400 {
						return errors.New("no connection to db")
					}
				}
				return nil
			}

			// wrap all requests in a transaction and set the length
			// of time doing things in the db to the log.
			tx := db.Begin()
			if tx.Error != nil {
				return errors.WithStack(tx.Error)
			}
			defer tx.Commit()

			c.Set("tx", tx)
			err := ef()
			if err != nil && errors.Cause(err) != errors.New("no connection to db") {
				tx.Rollback()
				return err
			}
			return nil
		}
	}
}

func paginate(c buffalo.Context, count int) *pop.Paginator {

	var (
		page     int
		per_page int
	)

	if p, err := strconv.Atoi(c.Param("page")); err == nil {
		page = p
	}

	if p, err := strconv.Atoi(c.Param("per_page")); err == nil {
		if p == 0 {
			p = 20
		}
		per_page = p
	}

	paginator := pop.NewPaginator(page, per_page)
	paginator.TotalEntriesSize = count
	paginator.TotalPages = paginator.TotalEntriesSize / paginator.PerPage
	if paginator.TotalEntriesSize%paginator.PerPage > 0 {
		paginator.TotalPages++
	}

	return paginator
}
