/**
 * Copyright 2015 @ to2.net.
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
