//go:build !linux

package ui

// flushStdin is a no-op on non-Linux platforms.
// macOS stty sane already handles terminal reset sufficiently.
func flushStdin() {}
