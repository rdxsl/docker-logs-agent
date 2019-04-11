package main

import (
	"github.com/rdxsl/docker-logs-agent/controllers"
	_ "github.com/rdxsl/docker-logs-agent/routers"

	"github.com/astaxie/beego"
)

var Version = "version notset"

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	controllers.Version = Version
	beego.Run()
}
