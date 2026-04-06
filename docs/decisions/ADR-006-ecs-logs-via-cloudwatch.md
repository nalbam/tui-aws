# ADR-006: ECS container logs via CloudWatch Logs instead of ECS API

## Status
Accepted

## Context
The ECS tab provides a deep-dive hierarchy: Clusters > Services > Tasks > Containers > Logs. Container logs need to be displayed in the TUI. There are two approaches:

1. **ECS API** — No direct log retrieval API exists. The `ecs:DescribeTasks` response includes `logDriver` and `logConfiguration` but not the actual log content.
2. **CloudWatch Logs API** — The `awslogs` log driver (default for Fargate and most ECS configurations) sends container logs to CloudWatch Logs with a predictable log group/stream naming convention.

## Decision
Use the CloudWatch Logs `GetLogEvents` API via the `CWL` client to fetch container logs. The log group and log stream are derived from the task definition's `logConfiguration`:

- Log group: from `logConfiguration.options["awslogs-group"]`
- Log stream: `{prefix}/{container-name}/{task-id}`

Implemented in `internal/aws/ecs.go:FetchContainerLogs()`, called from `internal/ui/tab_ecs/model.go`.

Relevant commit: `eb79479` (ECS deep dive — tasks, containers, logs, ECS Exec).

## Consequences

### Positive
- Works with the vast majority of ECS deployments (awslogs is the default and most common log driver)
- Reuses the existing `CWL` client already in the Clients struct (also used by CloudWatch tab)
- GetLogEvents supports `tailLines`-style retrieval with `Limit` parameter
- No additional AWS permissions beyond `logs:GetLogEvents`

### Negative
- Does not work with non-awslogs log drivers (fluentd, splunk, etc.) — shows "no logs available"
- Requires the task definition to have `logConfiguration` set (no logs for tasks without logging configured)

### Risks
- Log group naming conventions may vary across custom configurations (mitigated: extracted directly from task definition's logConfiguration, not hardcoded)
