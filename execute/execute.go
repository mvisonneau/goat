package execute

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

//CommandOut is a struct containing the return code, stdout, and stderr of an executed command
type CommandOut struct {
	Stdout string
	Stderr string
	Status int
}

//Command executes a given command + args and returns a CommandOut struct with an error if the command fails
func Command(command string, args []string, dir string) (CommandOut, error) {
	out := CommandOut{}
	cmd := exec.Cmd{}

	path, err := exec.LookPath(command)
	if err != nil {
		return out, err
	}

	cmd.Path = path
	cmd.Args = append([]string{path}, args...)

	if dir != "" {
		cmd.Dir = dir
	}

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Start(); err != nil {
		return out, fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				out.Stdout = cmdOut.String()
				out.Stderr = cmdErr.String()
				out.Status = status.ExitStatus()
				return out, fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			out.Stdout = cmdOut.String()
			out.Stderr = cmdErr.String()
			return out, fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	out.Stdout = cmdOut.String()
	out.Stderr = cmdErr.String()
	out.Status = 0
	return out, nil
}
