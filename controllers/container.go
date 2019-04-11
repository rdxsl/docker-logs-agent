package controllers

import (
	"strconv"

	"github.com/rdxsl/docker-logs-agent/models"

	"github.com/astaxie/beego"
)

// Operations about Container
type ContainerController struct {
	beego.Controller
}

// @Title Get
// @Description find object by objectid
// @Param	containerID		path 	string	true		"the containerID you want to get"
// @Param	tail		query 	string	false		"number of lines you want to get"
// @Success 200 {container} models.Container
// @Failure 403 :containerID is empty
// @router /:containerID/logs/?tail= [get]
func (o *ContainerController) Get() {
	containerID := o.Ctx.Input.Param(":containerID")
	tail := o.GetString("tail")

	if tail != "" {
		i, err := strconv.Atoi(tail)
		if err != nil || i < 0 {
			o.Data["json"] = map[string]string{"err": "container logs tail query string needs to be integer greater than 0"}
			o.ServeJSON()
			return
		}
	}

	if containerID != "" {
		ob, err := models.GetLog(containerID, tail)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = map[string]string{"containerID": ob}
		}
	}
	o.ServeJSON()
}
