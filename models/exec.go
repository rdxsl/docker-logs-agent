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
	ContainerID string `json:"ContainerID"`
	Result      string `json:"Result"`
}

func Exec(containerID string, Cmd1 ExecCmd) (execResult ExecResult, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(beego.AppConfig.String("docker_api_version")))
	if err != nil {
		return
	}
	c, err := cli.ContainerInspect(ctx, containerID)
	// need to handle error
	execResult.ContainerID = containerID

	cfg := &types.ExecConfig{
		Privileged:   true,
		Tty:          c.Config.Tty,
		Detach:       false,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}
	cfg.Cmd = Cmd1.Cmd
	response, err := cli.ContainerExecCreate(ctx, containerID, *cfg)

	if err != nil {
		err = fmt.Errorf("Error creating container exec: %v", err)
		fmt.Println(err)
		return
	}

	execID := response.ID
	if execID == "" {
		return
	}
	fmt.Println(execID)

	startCfg := types.ExecStartCheck{
		Detach: cfg.Detach,
		Tty:    cfg.Tty,
	}
	stream, err := cli.ContainerExecAttach(ctx, execID, startCfg)
	if err != nil {
		err = fmt.Errorf("Error attaching to container exec: %v", err)
		return
	}
	defer stream.Close()

	buf := new(bytes.Buffer)

	if c.Config.Tty {
		_, err = io.Copy(buf, stream.Reader)
	} else {
		_, err = stdcopy.StdCopy(buf, buf, stream.Reader)
	}
	buf.ReadFrom(stream.Reader)
	execResult.Result = buf.String()
	return
}
