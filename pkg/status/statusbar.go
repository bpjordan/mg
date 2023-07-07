package status

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type StatusBar struct {
    totalTasks uint
    remainingTasks []string
    wscol, wsrow uint16
    sigWinch, sigTerm chan os.Signal
    once sync.Once
}

func StartStatusBar(description string, tasks []string) (*StatusBar, error) {

    sb := &StatusBar{
	totalTasks: uint(len(tasks)),
	remainingTasks: tasks[:],
	sigWinch: make(chan os.Signal),
	sigTerm: make(chan os.Signal),
    }

    signal.Notify(sb.sigWinch, syscall.SIGWINCH)
    signal.Notify(sb.sigTerm, syscall.SIGTERM)

    sb.reserve()

    go func() { // Handle signals
	for {
	    select {
	    case <- sb.sigWinch:
		sb.reserve()
	    case <- sb.sigTerm:
		sb.Cleanup()
		os.Exit(255)
	    }
	}
    }()

    sb.render()
    return sb, nil
}

func (sb *StatusBar) PopTask(task string) error {
    for idx, v := range sb.remainingTasks {
	if task == v {
	    sb.remainingTasks = append(sb.remainingTasks[:idx], sb.remainingTasks[idx+1:]...)

	    sb.render()
	    return nil
	}
    }

    return fmt.Errorf("task %s not found", task)
}

func (sb *StatusBar) reserve() error {
    fmt.Printf("\x1B[0;%dr", sb.wsrow) // Drop existing margin reservation

    ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
    if err != nil {
        return err
    }

    sb.wsrow = ws.Row
    sb.wscol = ws.Col

    fmt.Print("\x1BD") // Return carriage
    fmt.Print("\x1B7") // Save cursor position
    fmt.Printf("\x1B[0;%dr", sb.wsrow-1) // Reserve bottom line
    fmt.Print("\x1B8") // Restore cursor position
    fmt.Print("\x1B[1A") // Move cursor up # lines

    return nil
}

func (sb *StatusBar) Cleanup() {
    sb.once.Do(func() {close(sb.sigTerm)})
    sb.once.Do(func() {close(sb.sigWinch)})

    fmt.Print("\x1B7")                 // Save the cursor position
    fmt.Printf("\x1B[0;%dr", sb.wsrow) // Drop margin reservation
    fmt.Printf("\x1B[%d;0f", sb.wsrow) // Move the cursor to the bottom line
    fmt.Print("\x1B[0K")               // Erase the entire line
    fmt.Print("\x1B8")                 // Restore the cursor position util new size is calculated

}

func (sb *StatusBar) render() {

    /// UPDATE POSITION
    fmt.Print("\x1B7") // save cursor position
    fmt.Print("\x1B[2K") // Erase current line
    fmt.Print("\x1B[0J") // Erase from cursor to end of screen
    fmt.Print("\x1B[?47h") // Save screen
    fmt.Print("\x1B[1J") // Erase from cursor to beginning of screen
    fmt.Print("\x1B[?47l") // Restore screen
    defer fmt.Print("\x1B8") // Restore cursor position

    numRemaining := sb.totalTasks - uint(len(sb.remainingTasks))

    /// Actually print stuff here
    line := fmt.Sprintf("[%d/%d] %s", numRemaining, sb.totalTasks, strings.Join(sb.remainingTasks, ", "))


    if len(line) > int(sb.wscol) {
	line = string([]rune(line)[:sb.wscol-1])
    }

    fmt.Printf("\x1B[%d;H", sb.wsrow)
    fmt.Print(line)

}
