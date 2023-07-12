package runtime

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/unix"
)

func (sb *ParallelRuntime) placeStatusBar() error {
    fmt.Printf("\x1B[0;%dr", sb.wsrow) // Drop existing margin reservation

    ws, err := unix.IoctlGetWinsize(unix.Stdout, unix.TIOCGWINSZ)
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

func (sb *ParallelRuntime) cleanupStatusBar() {

    fmt.Print("\x1B7")                 // Save the cursor position
    fmt.Printf("\x1B[0;%dr", sb.wsrow) // Drop margin reservation
    fmt.Printf("\x1B[%d;0f", sb.wsrow) // Move the cursor to the bottom line
    fmt.Print("\x1B[0K")               // Erase the entire line
    fmt.Print("\x1B8")                 // Restore the cursor position util new size is calculated

}

func (sb *ParallelRuntime) renderStatusBar() {

    /// UPDATE POSITION
    fmt.Print("\x1B7") // save cursor position
    fmt.Print("\x1B[2K") // Erase current line
    fmt.Print("\x1B[0J") // Erase from cursor to end of screen
    fmt.Print("\x1B[?47h") // Save screen
    fmt.Print("\x1B[1J") // Erase from cursor to beginning of screen
    fmt.Print("\x1B[?47l") // Restore screen
    defer fmt.Print("\x1B8") // Restore cursor position

    numRemaining := sb.totalTasks - uint(len(sb.activeTasks))

    /// Actually print stuff here
    number := color .New(color.FgYellow).Sprintf("[%d/%d]", numRemaining, sb.totalTasks)
    line := fmt.Sprintf("%s %s", number, strings.Join(sb.activeTasks, ", "))

    if len(line) > int(sb.wscol) {
	line = string([]rune(line)[:sb.wscol-1])
    }

    fmt.Printf("\x1B[%d;H", sb.wsrow) // Set cursor position to reserved row
    fmt.Print(line)

}
