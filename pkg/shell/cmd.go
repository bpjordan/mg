package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/bpjordan/multigit/pkg/manifest"
	"github.com/bpjordan/multigit/pkg/status"
	"github.com/fatih/color"
)

type shellResult struct {
	name string
	stdout, stderr string
	err error
	exit int
}

func RunParallelCmd(bin string, args []string, man manifest.Manifest) (numSuccess, numFailed, numError uint) {

	taskFinished := make(chan shellResult)
	taskStarted := make(chan string)

	defer close(taskFinished)

	for _, repo := range man.Repos {
		go startCmd(bin, args, repo.Name, repo.Path, taskStarted, taskFinished)
	}

	sb, err := status.StartStatusBar(uint(len(man.Repos)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create status bar: \n", err.Error())
	}
	defer sb.Cleanup()

	for {
		select {
		case <- sb.Finished():
			sb.Cleanup()
			return

		case task := <- taskStarted:
			sb.PushTask(task)

		case result := <- taskFinished:
			printTaskReport(result, &numSuccess, &numFailed, &numError)
			sb.PopTask(result.name)
		}

	}
}

func startCmd(
	cmd string, args []string, name, dir string,
	start chan string, finish chan shellResult,
) {

	var stdout, stderr bytes.Buffer

	result := shellResult{}
	result.name = name

	task := exec.Command(cmd, args...)
	task.Dir = dir
	task.Stdout = &stdout
	task.Stderr = &stderr

	start <- name
	err := task.Run()

	result.stdout = stdout.String()
	result.stderr = stderr.String()

	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		result.exit = exitErr.ExitCode()
	} else {
		result.err = err
	}

	finish <- result
}

func printTaskReport(result shellResult, numSuccess, numFailed, numError *uint) {
	switch {
	case result.err != nil:
		*numError++
		fmt.Printf("Error executing command in %s: %e\n", result.name, result.err)

	case result.exit != 0:
		*numFailed++
		fmt.Printf("%s (repo %s) (return code %d):\n", color.RedString("Error"), result.name, result.exit)
		fmt.Println(result.stderr)

	default:
		*numSuccess++
		fmt.Printf("%s (%s)\n", color.GreenString("Success"), result.name)

		if result.stdout != "" {
			fmt.Println(result.stdout)
		}

		if result.stderr != "" {
			fmt.Println(color.HiYellowString("Warning"), "(stderr):")
			fmt.Println(result.stderr)
		}

	}
}
