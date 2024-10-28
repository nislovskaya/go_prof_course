package main

import (
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return int(syscall.EINVAL)
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	var envVars []string
	for key, value := range env {
		if value.NeedRemove {
			envVars = append(envVars, key+"=")
		} else {
			envVars = append(envVars, key+"="+value.Value)
		}
	}

	command.Env = append(os.Environ(), envVars...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		returnCode = 1
	} else {
		returnCode = 0
	}

	return returnCode
}
