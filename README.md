# Gof   
The development framework with golang.
**Gof lets you write web/server apps in Golang.**

## Web Framework ##
Include routes,server,mvc,paging,interceptor,template..

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
        		ctx.App.Template().Execute(ctx.ResponseWriter,
        		gof.TemplateDataMap{
        			"变量名": "变量值",
        			"SysName":sysName,
        		},"views/index.html")
        	})

        	var in = getInterceptor(app,routes)
        	go http.ListenAndServe(":8080",in)

        	log.Println("[ OK] - web is listening on port :8080.")
        	var ch = make(chan int,1)
        	<- ch

Details in https://github.com/atnet/gof/blob/master/web/demo/demo.go .


