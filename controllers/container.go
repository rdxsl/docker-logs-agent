package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/rdxsl/docker-agent-proxy/models"
)

// Operations about Container
type ContainerController struct {
	beego.Controller
}

// @Title Logs
// @Description return logs by containerID
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
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// @Title Exec Post
// @Description Exec a container
// @Param	containerID		path 	string	true		"the containerID you want to run exec"
// @Param	body		body 	models.ExecCmd	true		"The object content"
// @Success 200 {container} models.Container
// @Failure 400 :containerID is empty or bad Exec Cmd Json
// @Failure 409 :ExecCmd JSON can't be Unmarshaled
// @router /:containerID/exec [post]
func (o *ContainerController) Exec() {
	containerID := o.Ctx.Input.Param(":containerID")

	if containerID == "" {
		o.Ctx.Output.SetStatus(400)
		o.Data["json"] = map[string]string{"Error": "containerID is empty"}
		o.ServeJSON()
		return
	}

	var execCmd models.ExecCmd
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &execCmd)
	if err != nil {
		fmt.Println(o.Ctx.Input.RequestBody)
		o.Ctx.Output.SetStatus(400)
		o.Data["json"] = map[string]string{"Error": "Exec Cmd Json format wrong. " + err.Error()}
		o.ServeJSON()
		return
	}

	fmt.Println(containerID)
	execResult, err := models.Exec(containerID, execCmd)
	if err != nil {
		o.Ctx.Output.SetStatus(400)
		o.Data["json"] = map[string]string{"Error": "Exec Cmd Error. " + err.Error()}
		o.ServeJSON()
		return
	}
	o.Data["json"] = execResult
	o.ServeJSON()
	return
}
