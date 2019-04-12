package models

import (
	"bytes"

	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// var (
// 	Containers map[string]*Container
// )
//
// type Container struct {
// 	containerID string
// 	Logs        []byte
// }
//
// func init() {
// }
type Logs struct {
	ContainerLog       containerLog       `json:"containerLog"`
	Base64ContainerLog base64ContainerLog `json:"base64ContainerLog"`
}

type containerLog struct {
	ContainerID string `json:"containerID"`
	Logs        string `json:"Logs"`
}

type base64ContainerLog struct {
	ContainerID string `json:"base64containerID"`
	Base64Logs  []byte `json:"base64Logs"`
}

func GetLog(containerID string, tail string) (Logs, error) {
	return dockerContainerLogs(containerID, tail)
}

func dockerContainerLogs(containerID string, tail string) (Logs, error) {
	var l Logs
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(beego.AppConfig.String("docker_api_version")))
	if err != nil {
		return l, err
	}

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		Tail:       tail,
	}
	// Replace this ID with a container that really exists
	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return l, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	s := buf.Bytes()

	l.Base64ContainerLog = base64ContainerLog{containerID, s}
	l.ContainerLog = containerLog{containerID, buf.String()}
	return l, nil

}
