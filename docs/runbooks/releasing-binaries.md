# Runbook: Releasing Binaries

## Purpose
Build and distribute cross-compiled tui-aws binaries for all supported platforms.

## Prerequisites
- [ ] Go 1.25+ installed (`go version`)
- [ ] All tests pass (`make test`)
- [ ] `go vet ./...` reports no issues
- [ ] Version number updated in `Makefile` (`VERSION := X.Y.Z`)

## Procedure

### Step 1: Update version
Edit `Makefile` and set the new version:
```bash
# In Makefile, update:
VERSION := 0.2.0
```

### Step 2: Run tests and static analysis
```bash
make test && go vet ./...
```
Expected: all tests pass, no vet warnings.

### Step 3: Cross-compile
```bash
make clean && make build-all
```
Expected output in `dist/`:
```
dist/tui-aws-linux-amd64
dist/tui-aws-linux-arm64
dist/tui-aws-darwin-arm64
dist/tui-aws-darwin-amd64
```

### Step 4: Verify binaries
```bash
file dist/tui-aws-*
```
Expected: each binary shows correct architecture (ELF for Linux, Mach-O for macOS, amd64/arm64).

### Step 5: Test local binary
```bash
make build && ./tui-aws
```
Verify: TUI launches, tabs switch, profile/region selectors work.

### Step 6: Tag the release
```bash
git add Makefile
git commit -m "release: v0.2.0"
git tag -a v0.2.0 -m "v0.2.0"
git push origin main --tags
```

### Step 7: Create GitHub release (optional)
```bash
gh release create v0.2.0 dist/tui-aws-* --title "v0.2.0" --notes "Release notes here"
```

## Verification
- Each binary runs on its target platform
- `./tui-aws` shows the correct version (if version flag is implemented)
- GitHub release page lists all 4 binaries

## Rollback
```bash
# Delete tag locally and remotely
git tag -d v0.2.0
git push origin --delete v0.2.0

# Revert version commit if needed
git revert HEAD
```

## Related
- Architecture: docs/architecture.md#infrastructure
- Makefile: `build-all` target
