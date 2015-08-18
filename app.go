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
	"github.com/jrsix/gof/db"
	"github.com/jrsix/gof/log"
)

// 应用当前的上下文
var CurrentApp App

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
