/**
 * Copyright 2015 @ S1N1 Team.
 * name : doc.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package web

// Http Context Session
//
//    func (this *mainC) Index(ctx *web.Context) {
//        var t int64
//        v := ctx.Session().Get("num")
//        if i==1 || v == nil{
//            t = time.Now().Unix()
//            ctx.Session().Set("num",t)
//            ctx.Session().Save()
//            i = 2
//        }else{
//            t = v.(int64)
//        }
//
//        ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter, func(m *map[string]interface{}) {
//            (*m)["unix"] = t
//        },
//        "views/main/index.html",
//        "views/main/inc/header.html",
//        "views/main/inc/footer.html")
//    }
//
//  Destroy Session
//  ctx.Session().Destroy()
//
