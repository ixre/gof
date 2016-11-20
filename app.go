/**
 * Copyright 2014 @ z3q.net.
 * name : app1.go
 * author : jarryliu
 * date : 2015-04-27 20:43:
 * description :
 * history :
 */
package gof

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/storage"
)

// 应用当前的上下文
var CurrentApp App

type App interface {
	// Provided db access
	Db() db.Connector
	//Orm() orm.Orm

	// Return a Wrapper for GoLang template.
	Template() *Template
	// Return application configs.
	Config() *Config
	// Storage
	Storage() storage.Interface
	// Return a logger
	Log() log.ILogger
	// Application is running debug mode
	Debug() bool
}
