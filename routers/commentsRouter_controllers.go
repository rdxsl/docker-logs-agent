package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/rdxsl/docker-logs-agent/controllers:ContainerController"] = append(beego.GlobalControllerRouter["github.com/rdxsl/docker-logs-agent/controllers:ContainerController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/logs/:containerId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
