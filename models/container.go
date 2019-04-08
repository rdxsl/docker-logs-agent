package models

import (
	"bytes"

	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var (
	Containers map[string]*Container
)

type Container struct {
	ContainerId string
	Logs        string
}

func init() {
	Containers = make(map[string]*Container)
}

func GetLog(ContainerId string) (string, error) {
	return dockerContainerLogs(ContainerId)
}

func dockerContainerLogs(ContainerId string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(beego.AppConfig.String("docker_api_version")))
	if err != nil {
		return "can't connect to docker api", err
	}

	options := types.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := cli.ContainerLogs(ctx, ContainerId, options)
	if err != nil {
		return "no such container", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	s := buf.String()
	return s, nil

}
