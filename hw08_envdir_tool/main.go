package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Error reading env dir: %v\n", err)
		os.Exit(1)
	}

	commandsToRun := append([]string{command}, args...)

	returnCode := RunCmd(commandsToRun, env)

	os.Exit(returnCode)
}
