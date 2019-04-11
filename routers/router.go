// @APIVersion 1.0.0
// @Title docker-logs-agent
// @Description Web wrapper to access docker logs via unix socket `/var/run/docker.sock`. Only use this in a secure network environment.
// @Contact jxie@riotgames.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/rdxsl/docker-logs-agent/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/containers",
			beego.NSInclude(
				&controllers.ContainerController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
