/**
 * Copyright 2015 @ z3q.net.
 * name : controller.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package mvc

type Controller interface{}

// Generate controller instance
type ControllerGenerate func() Controller
