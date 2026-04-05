# ADR-004: AWS profile credential fallback chain

## Status
Accepted

## Context
Users run tui-aws in diverse environments: local dev with named profiles, EC2 instances with instance roles, containers with task roles, and SSO configurations. The tool needs to gracefully handle missing or invalid profiles without crashing or showing cryptic SDK errors.

## Decision
Implement a multi-step fallback chain for profile resolution:
1. Use saved profile from `~/.tui-aws/config.json` if valid
2. Try `default` profile from `~/.aws/credentials`
3. Fallback to first named profile found in credentials/config
4. Final fallback to instance role (no `--profile` flag)

Each step validates credentials via `STS GetCallerIdentity` before accepting.

## Consequences

### Positive
- Works out-of-the-box on EC2 instances, ECS tasks, and local dev
- Invalid saved config doesn't block startup
- User sees which profile was selected in the status bar

### Negative
- Startup may be slower when multiple validation calls fail before finding a valid profile
- The "first named profile" heuristic may not pick the user's preferred profile

### Risks
- STS calls during startup add latency on slow networks (mitigated: fail-fast with short timeout)
