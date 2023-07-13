package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ParallelRuntime struct {
    totalTasks, remainingTasks uint
    activeTasks []string
    wscol, wsrow uint16
    sigWinch, sigTerm chan os.Signal
    finished chan struct{}
    once sync.Once
    ctx context.Context
    cancel context.CancelFunc
}

func Start(ctx context.Context, totalTasks uint) (*ParallelRuntime, error) {

    ctx, cancel := context.WithCancel(ctx)
    sb := &ParallelRuntime{
	totalTasks: totalTasks,
	remainingTasks: totalTasks,
	activeTasks: make([]string, 0, totalTasks),
	sigWinch: make(chan os.Signal),
	sigTerm: make(chan os.Signal),
	ctx: ctx,
	cancel: cancel,
    }

    signal.Notify(sb.sigWinch, syscall.SIGWINCH)
    signal.Notify(sb.sigTerm, syscall.SIGTERM)

    sb.placeStatusBar()

    go func() { // Handle signals
	for {
	    select {
	    case <- sb.sigWinch:
		sb.placeStatusBar()
	    case <- sb.sigTerm:
		sb.cancel()
	    }
	}
    }()

    sb.renderStatusBar()
    return sb, nil
}

func (rt *ParallelRuntime) Cancel() {
    rt.cancel()
}

func (sb *ParallelRuntime) Cleanup() {
    sb.cleanupStatusBar()
}

func (rt *ParallelRuntime) Context() context.Context {
    return rt.ctx
}

func (sb *ParallelRuntime) PushTask(task string) {
    sb.activeTasks = append(sb.activeTasks, task)
    sb.renderStatusBar()
}

func (sb *ParallelRuntime) PopTask(task string) error {
    for idx, v := range sb.activeTasks {
	if task == v {
	    sb.activeTasks = append(sb.activeTasks[:idx], sb.activeTasks[idx+1:]...)
	    sb.remainingTasks--
	    sb.renderStatusBar()

	    if sb.remainingTasks == 0 {
		sb.cancel()
	    }

	    return nil
	}
    }

    return fmt.Errorf("task %s not found", task)
}

func (sb *ParallelRuntime) Finished() <-chan struct{} {
    return sb.ctx.Done()
}

