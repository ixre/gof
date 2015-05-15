# Gof   
Golang micro common framework

## Web Framework ##
Include routes,server,mvc,paging,interceptor

example:

        app := &HttpApp{}
        	routes := &web.RouteMap{}

        	routes.Add("/[0-9]/*",func(ctx *web.Context){
        		ctx.ResponseWriter.Write([]byte("数字路径"))
        	})

        	routes.Add("/[a-z]$",func(ctx *web.Context){
        		ctx.ResponseWriter.Write([]byte("字母路径"))
        	})

        	routes.Add("/",func(ctx *web.Context){
        		sysName := ctx.App.Config().GetString("SYS_NAME")
        		ctx.ResponseWriter.Write([]byte("Hello,Gof with "+ sysName+"."))
        		ctx.ResponseWriter.Header().Set("Content-Type","text/html")
        		return
        		ctx.App.Template().Execute(ctx.ResponseWriter,
        		func(v *map[string]interface{}){
        			(*v)["变量名"] = "变量值"
        		},"views/index.html")
        	})

        	var in = getInterceptor(app,routes)
        	go http.ListenAndServe(":8080",in)

        	log.Println("[ OK] - web is listening on port :8080.")
        	var ch = make(chan int,1)
        	<- ch

Details in https://github.com/atnet/gof/blob/master/web/demo/web_demo.go .


