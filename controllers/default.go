package controllers

import (
	"github.com/astaxie/beego"
)

var Version string

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	VersionJson := map[string]string{"docker-logs-agent": Version}
	VersionJson["docker-api"] = beego.AppConfig.String("docker_api_version")
	this.Data["json"] = VersionJson
	this.ServeJSON()
}
