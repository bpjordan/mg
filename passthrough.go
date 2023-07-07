package main

import (
	"fmt"
	"os/exec"
)

type shellResult struct {
	dir string
	stdout, stderr string
	err error
	exit int
}

func RunCmd(bin string, args, repo_paths []string) {

	results := make(chan shellResult)

	for _, path := range repo_paths {
		go startCmd("git", args, path, results)
	}

	for i := 0; i < len(repo_paths); i++ {
		result := <- results

		if result.err != nil {
			fmt.Printf("Error executing command in %s: %e\n", result.dir, result.err)
		} else if result.exit != 0 {
			fmt.Printf("Error (repo %s) (return code %d):\n", result.dir, result.exit)
			fmt.Println(result.stderr)
		} else {
			fmt.Printf("Success (repo %s)", result.dir)
			fmt.Println(result.stdout)

			if len(result.stderr) > 0 {
				fmt.Println("Warning:")
				fmt.Println(result.stderr)
			}
			fmt.Println()
		}
	}
}

func startCmd(cmd string, args []string, dir string, resultPipe chan shellResult) {

	result := shellResult{}
	result.dir = dir

	task := exec.Command(cmd, args...)
	task.Dir = dir

	stdout, err := task.StdoutPipe()
	if err != nil {
		result.err = err
		resultPipe <- result
		return
	}

	stderr, err := task.StderrPipe()
	if err != nil {
		result.err = err
		resultPipe <- result
		return
	}

	err = task.Run()
	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		result.exit = exitErr.ExitCode()
	} else {
		result.err = err
	}

	buf := make([]byte, 0)
	_, err = stdout.Read(buf)
	result.stdout = string(buf)

	_, err = stderr.Read(buf)
	result.stderr = string(buf)

	resultPipe <- result
}
