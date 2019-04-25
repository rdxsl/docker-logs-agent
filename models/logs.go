package models

import (
	"bytes"
	"io"

	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

type Logs struct {
	ContainerLog containerLog `json:"containerLog"`
}

type containerLog struct {
	ContainerID string `json:"containerID"`
	Logs        string `json:"Logs"`
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
		ShowStderr: true,
		Tail:       tail,
	}
	c, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return l, err
	}

	// Replace this ID with a container that really exists
	reader, err := cli.ContainerLogs(ctx, c.ID, options)
	if err != nil {
		return l, err
	}

	defer reader.Close()

	buf := new(bytes.Buffer)

	if c.Config.Tty {
		_, err = io.Copy(buf, reader)
	} else {
		_, err = stdcopy.StdCopy(buf, buf, reader)
	}

	buf.ReadFrom(reader)

	l.ContainerLog = containerLog{containerID, buf.String()}
	return l, nil

}
