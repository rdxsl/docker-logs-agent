package controllers

import (
	"github.com/astaxie/beego"
)

var Version string

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	VersionJson := map[string]string{"docker-agent-proxy": Version}
	VersionJson["docker-api"] = beego.AppConfig.String("dockerApiVersion")
	this.Data["json"] = VersionJson
	this.ServeJSON()
}
