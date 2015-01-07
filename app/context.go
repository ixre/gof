package app

import (
	"github.com/newmin/gof"
	"github.com/newmin/gof/db"
	"github.com/newmin/gof/log"
	"github.com/newmin/gof/web"
)

type Context interface {
	// Provided db access
	Db() db.Connector

	// Return a Wrapper for golang template.

	Template() *web.TemplateWrapper

	// Return application configs.
	Config() *gof.Config

	// Return a logger
	Log() log.ILogger

	// Get a reference of AppContext
	Source() interface{}

	// Application is running debug mode
	Debug() bool
}
