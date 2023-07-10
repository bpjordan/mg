package shell

import (
	"fmt"
	"io"
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

	results := make(chan shellResult)
	defer close(results)

	for _, repo := range man.Repos {
		go startCmd("git", args, repo.Name, repo.Path, results)
	}

	sb, err := status.StartStatusBar(man.Names())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create status bar: \n", err.Error())
	}
	defer sb.Cleanup()

	for i := 0; i < len(man.Repos); i++ {
		result := <- results

		switch {
		case result.err != nil:
			numError++
			fmt.Printf("Error executing command in %s: %e\n", result.name, result.err)

		case result.exit != 0:
			numFailed++
			fmt.Printf("%s (repo %s) (return code %d):\n", color.RedString("Error"), result.name, result.exit)
			fmt.Println(result.stderr)

		default:
			numSuccess++
			fmt.Printf("%s (%s)\n", color.GreenString("Success"), result.name)

			if result.stdout != "" {
				fmt.Println(result.stdout)
			}

			if result.stderr != "" {
				fmt.Println(color.HiYellowString("Warning"), "(stderr):")
				fmt.Println(result.stderr)
			}

		}

		sb.PopTask(result.name)
	}

	return
}

func startCmd(cmd string, args []string, name, dir string, resultPipe chan shellResult) {

	result := shellResult{}
	result.name = name

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

	err = task.Start()

	stderrBytes, err := io.ReadAll(stderr)
	result.stderr = string(stderrBytes)
	if err != nil {
		result.err = err
		resultPipe <- result
		return
	}

	stdoutBytes, err := io.ReadAll(stdout)
	result.stdout = string(stdoutBytes)
	if err != nil {
		result.err = err
		resultPipe <- result
		return
	}

	err = task.Wait()
	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		result.exit = exitErr.ExitCode()
	} else {
		result.err = err
	}


	resultPipe <- result
}
