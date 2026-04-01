---
name: sync-docs
description: Synchronize project documentation with current code state. Use when asked to update docs, audit documentation, or sync CLAUDE.md/steering files.
---

# Sync Docs Skill

## Actions
1. **Steering rules audit** — Check `.kiro/steering/` files match current code patterns
2. **Architecture doc sync** — Update `docs/architecture.md` to reflect current system
3. **ADR audit** — Check recent commits, suggest new ADRs for undocumented decisions
4. **README sync** — Update project structure section to match actual layout
5. **Report** — List all changes made
