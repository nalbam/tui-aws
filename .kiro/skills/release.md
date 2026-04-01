---
name: release
description: Automate the release process with validation checks. Use when asked to create a release, bump version, or tag.
---

# Release Skill

## Procedure
1. **Pre-release** — Clean tree (`git status`), all tests pass (`make test`)
2. **Version** — Review changes since last tag, apply semver (MAJOR/MINOR/PATCH)
3. **Changelog** — Group by Added/Changed/Fixed/Removed
4. **Build** — Update Makefile VERSION, `make build-all`, create git tag
5. **Summary** — Display version bump, key changes, next steps
