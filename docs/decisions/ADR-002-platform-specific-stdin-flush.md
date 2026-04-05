# ADR-002: Platform-specific stdin flush with build tags

## Status
Accepted

## Context
After SSM/ECS Exec sessions, the terminal stdin buffer contains leftover bytes from the interactive session. Without flushing, these bytes are interpreted as Bubble Tea keystrokes, causing erratic behavior when the TUI resumes. The `unix.TCFLSH` syscall is Linux-specific and doesn't compile on macOS.

## Decision
Split the stdin flush into platform-specific files using Go build tags:
- `flush_linux.go` — uses `unix.TCFLSH` (TCIFLUSH) for real buffer flush
- `flush_other.go` — no-op stub for macOS and other platforms

Both files are in `internal/ui/` and export the same function signature.

## Consequences

### Positive
- Cross-platform compilation works (linux + darwin, amd64 + arm64)
- Linux gets proper stdin flush (most common deployment target)
- Clean separation via build tags instead of runtime detection

### Negative
- macOS users may see occasional stale keystrokes after SSM sessions
- Adding new platform-specific behavior requires new build tag files

### Risks
- None significant — macOS is primarily for development, Linux for production use
