# Build and Deploy

Build cross-platform binaries for tui-aws.

## Steps

1. Pre-build checks:
   ```bash
   go vet ./...
   go test ./... -v
   ```

2. If checks pass, build all platforms:
   ```bash
   make build-all
   ```

3. Verify build artifacts:
   ```bash
   ls -la dist/
   file dist/tui-aws-*
   ```

4. Report:
   - Binary sizes for each platform
   - Build status (success/failure)
   - Version from Makefile

## Notes
- Requires Go 1.25+ for cross-compilation
- Output goes to `dist/` directory
- Platforms: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
