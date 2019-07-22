package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/rdxsl/docker-agent-proxy/controllers:ContainerController"] = append(beego.GlobalControllerRouter["github.com/rdxsl/docker-agent-proxy/controllers:ContainerController"],
        beego.ControllerComments{
            Method: "Exec",
            Router: `/:containerID/exec`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/rdxsl/docker-agent-proxy/controllers:ContainerController"] = append(beego.GlobalControllerRouter["github.com/rdxsl/docker-agent-proxy/controllers:ContainerController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:containerID/logs/?tail=`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
