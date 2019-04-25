package models

import (
	"bytes"
	"fmt"
	"io"

	"github.com/astaxie/beego"
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

type execOptions struct {
	detachKeys  string
	interactive bool
	tty         bool
	detach      bool
	user        string
	privileged  bool
	env         opts.ListOpts
	workdir     string
	container   string
	command     []string
}

func newExecOptions() execOptions {
	return execOptions{env: opts.NewListOpts(opts.ValidateEnv)}
}

func LogTest(containerID string) (err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(beego.AppConfig.String("docker_api_version")))
	if err != nil {
		return err
	}
	c, err := cli.ContainerInspect(ctx, containerID)
	fmt.Println("LogTest" + containerID)

	cfg := &types.ExecConfig{
		Privileged:   true,
		Tty:          c.Config.Tty,
		Detach:       false,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}
	cfg.Cmd = []string{"pwd"}
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
	fmt.Println(buf.String())
	return
}
