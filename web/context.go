/**
 * Copyright 2015 @ S1N1 Team.
 * name : context.go
 * author : newmin
 * date : -- :
 * description :
 * history :
 */

package web
import(
    "net/http"
    "github.com/atnet/gof"
)


type HttpContextFunc func(*Context)

type Context struct {
    App gof.App
    ResponseWriter http.ResponseWriter
    Request *http.Request
}

func NewContext(app gof.App,rsp http.ResponseWriter,req *http.Request)*Context{
    return &Context{
        App:app,
        ResponseWriter:rsp,
        Request:req,
    }
}
