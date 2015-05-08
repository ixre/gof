/**
 * Copyright 2015 @ S1N1 Team.
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
