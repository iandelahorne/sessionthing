package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/lflux/sessionthing/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
