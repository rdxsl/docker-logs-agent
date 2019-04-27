package models

import (
	"bytes"
	"fmt"
	"io"

	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

type ExecCmd struct {
	Cmd []string `json:"Cmd"`
}

type ExecResult struct {
	ContainerID string  `json:"ContainerID"`
	ExecCmd     ExecCmd `json:"ExecCmd"`
	ExecResult  string  `json:"ExecResult"`
}

func Exec(containerID string, Cmd ExecCmd) (execResult ExecResult, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(beego.AppConfig.String("dockerApiVersion")))
	if err != nil {
		return
	}
	c, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		err = fmt.Errorf("Error in container exec ContainerInspect: %v", err)
		return
	}
	execResult.ContainerID = containerID

	cfg := &types.ExecConfig{
		Privileged:   true,
		Tty:          c.Config.Tty,
		Detach:       false,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}

	cfg.Cmd = Cmd.Cmd
	response, err := cli.ContainerExecCreate(ctx, containerID, *cfg)

	if err != nil {
		err = fmt.Errorf("Error in container exec ContainerExecCreate: %v", err)
		return
	}

	execID := response.ID
	if execID == "" {
		err = fmt.Errorf("Error in container exec execID is empty")
		return
	}

	startCfg := types.ExecStartCheck{
		Detach: cfg.Detach,
		Tty:    cfg.Tty,
	}

	stream, err := cli.ContainerExecAttach(ctx, execID, startCfg)
	if err != nil {
		err = fmt.Errorf("Error attaching to container exec ContainerExecAttach: %v", err)
		return
	}
	defer stream.Close()

	buf := new(bytes.Buffer)

	if c.Config.Tty {
		_, err = io.Copy(buf, stream.Reader)
		if err != nil {
			err = fmt.Errorf("Error attaching to container exec io.Copy: %v", err)
			return
		}
	} else {
		_, err = stdcopy.StdCopy(buf, buf, stream.Reader)
		if err != nil {
			err = fmt.Errorf("Error attaching to container exec stdcopy.StdCopy: %v", err)
			return
		}
	}

	buf.ReadFrom(stream.Reader)
	execResult.ExecResult = buf.String()
	execResult.ExecCmd = Cmd
	return
}
