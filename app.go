/**
 * Copyright 2014 @ S1N1 Team.
 * name : app1.go
 * author : jarryliu
 * date : 2015-04-27 20:43:
 * description :
 * history :
 */
package gof

import (
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
)

type App interface {
	// Provided db access
	Db() db.Connector
	// Return a Wrapper for GoLang template.
	Template() *Template
	// Return application configs.
	Config() *Config
	// Storage
	Storage() Storage
	// Return a logger
	Log() log.ILogger
	// Application is running debug mode
	Debug() bool
}
