package shell

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/bpjordan/multigit/pkg/manifest"
	"github.com/bpjordan/multigit/pkg/runtime"
	"github.com/fatih/color"
)

type shellResult struct {
	name string
	stdout, stderr string
	err error
	exit int
}

func RunParallelCmd(rt *runtime.ParallelRuntime, bin string, args []string, man manifest.Manifest, verbose int) (numSuccess, numFailed, numError uint) {

	taskFinished := make(chan shellResult)
	taskStarted := make(chan string)

	defer close(taskFinished)

	for _, repo := range man.Repos() {
		go startCmd(
			rt, bin, args,
			repo.Name, repo.Path,
			taskStarted, taskFinished,
		)
	}

	for {
		select {
		case <- rt.Finished():
			rt.Cleanup()
			return

		case task := <- taskStarted:
			rt.PushTask(task)

		case result := <- taskFinished:
			printTaskReport(result, &numSuccess, &numFailed, &numError, verbose)
			rt.PopTask(result.name)
		}

	}
}

func startCmd(
	rt *runtime.ParallelRuntime,
	cmd string, args []string, name, dir string,
	start chan string, finish chan shellResult,
) {

	var stdout, stderr bytes.Buffer

	result := shellResult{}
	result.name = name

	task := exec.CommandContext(rt.Context(), cmd, args...)
	task.Dir = dir
	task.Stdout = &stdout
	task.Stderr = &stderr

	if rt.Acquire() != nil {
		return
	}

	start <- name
	err := task.Run()
	rt.Release()

	result.stdout = stdout.String()
	result.stderr = stderr.String()

	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		result.exit = exitErr.ExitCode()
	} else {
		result.err = err
	}

	finish <- result
}

func printTaskReport(result shellResult, numSuccess, numFailed, numError *uint, verbose int) {
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

		if result.stdout != "" && verbose > 0{
			fmt.Println(result.stdout)
		}

		if result.stderr != "" {
			fmt.Printf("%s (%s)\n", color.HiYellowString("Warning"), result.name)
			fmt.Println(result.stderr)
		}

	}
}
