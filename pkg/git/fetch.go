package git

import (
	"fmt"

	"github.com/bpjordan/mg/pkg/manifest"
	"github.com/bpjordan/mg/pkg/runtime"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

type errorReport struct {
    Name string
    Err error
}

type taskReport struct {
    Name string
    Updated bool
    Err error
}

type FetchReport struct {
    // Tasks which resulted in updates to a repo
    Updated uint
    // Tasks which finished successfully but did not update a repo
    NoChange uint
    // Tasks which were successfully started but did not succeed
    Failed uint
    // Tasks which failed to begin
    Error uint
}

func Fetch(rt *runtime.ParallelRuntime, manifest manifest.Manifest) (*FetchReport, error) {

    taskStarted := make(chan string)
    taskFinished := make(chan taskReport)
    taskError := make(chan errorReport)

    for _, repo := range manifest.Repos() {
        go fetchRepo(rt, repo, taskStarted, taskFinished, taskError)
    }

    var report FetchReport

    verbose := viper.GetInt("verbose")
    for {

        select {
        case <- rt.Finished():
            return &report, nil
        case task := <- taskStarted:
            rt.PushTask(task)
        case task := <- taskFinished:
            rt.PopTask(task.Name)
            printTaskReport(task, verbose)
            switch {
            case task.Err != nil:
                report.Failed++
            case task.Updated:
                report.Updated++
            default:
                report.NoChange++
            }
        case task := <- taskError:
            rt.DecrementCounter()
            fmt.Printf("Error starting task for %s: %s", task.Name, task.Err.Error())
            report.Error++
        }
    }
}

func printTaskReport(report taskReport, verbose int) {
    if report.Err != nil {
        fmt.Printf("%s (%s)\n", color.RedString("ERROR"), report.Name)
        println(report.Err.Error())
    } else if report.Updated {
        fmt.Printf("%s (%s)\n", color.HiGreenString("UPDATED"), report.Name)
    } else if verbose > 0 {
        fmt.Printf("%s (%s)\n", color.GreenString("NO CHANGES"), report.Name)
    }
}

func fetchRepo(
    rt *runtime.ParallelRuntime,
    repoMeta manifest.Repository,
    taskStarted chan<- string, taskFinished chan<- taskReport, taskError chan<- errorReport,
) {
    r, err := git.PlainOpen(repoMeta.Path)
    if err != nil {
        taskError <- errorReport{
            repoMeta.Name,
            err,
        }
    }

    err = rt.Acquire()
    if err != nil {
        taskError <- errorReport{
            repoMeta.Name,
            fmt.Errorf("Failed to acquire task lock"),
        }
    }

    taskStarted <- repoMeta.Name

    err = r.FetchContext(rt.Context(), &git.FetchOptions{})

    rt.Release()

    switch err {
    case nil:
        taskFinished <- taskReport{
            repoMeta.Name,
            true,
            nil,
        }
    case git.NoErrAlreadyUpToDate:
        taskFinished <- taskReport{
            repoMeta.Name,
            false,
            nil,
        }
    default:
        taskFinished <- taskReport{
            repoMeta.Name,
            false,
            err,
        }
    }
}

