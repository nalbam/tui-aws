//go:build linux

package ui

import (
	"os"

	"golang.org/x/sys/unix"
)

// flushStdin discards any residual bytes in the stdin input buffer.
func flushStdin() {
	unix.IoctlSetInt(int(os.Stdin.Fd()), unix.TCFLSH, unix.TCIFLUSH) //nolint:errcheck
}
