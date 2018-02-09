package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

func FlashHandler(c buffalo.Context) error {
	c.Flash().Add("success", "Success flash")
	return c.Redirect(http.StatusFound, "/")
}
