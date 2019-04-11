package controllers

import (
	"github.com/rdxsl/docker-logs-agent/models"

	"github.com/astaxie/beego"
)

// Operations about Container
type ContainerController struct {
	beego.Controller
}

// @Title Get
// @Description find object by objectid
// @Param	containerId		path 	string	true		"the containerId you want to get"
// @Success 200 {container} models.Container
// @Failure 403 :containerID is empty
// @router /:containerId/logs [get]
func (o *ContainerController) Get() {
	containerId := o.Ctx.Input.Param(":containerId")
	if containerId != "" {
		ob, err := models.GetLog(containerId)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = map[string]string{"containerId": ob}
		}
	}
	o.ServeJSON()
}
