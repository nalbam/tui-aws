# ECS Tab
## Role
ECS deep dive — Clusters → Services → Tasks → Containers → Logs → ECS Exec.
## Key Files
- `model.go` — 15-state viewState for drill-down hierarchy, ECSExecRequest/ECSExecDoneMsg messages
- `table.go` — Per-level table columns (cluster, service, task, container)
- `detail.go` — Cluster/service/task/container/task-def detail + CloudWatch log viewer
## Drill-Down Hierarchy
```
ClusterList → ClusterAction → ServiceList → ServiceAction → TaskList → TaskAction → ContainerList → ContainerAction
                  ↓                ↓              ↓                ↓                     ↓
             ClusterDetail    ServiceDetail    TaskDetail     TaskDefDetail          ContainerDetail / Logs / ECS Exec
```
## Rules
- ECS Exec sends `ECSExecRequest` to RootModel (intercepted like SSM flow via `tea.Exec`)
- `ECSExecDoneMsg` is forwarded back from RootModel after exec session ends
- Container logs fetched from CloudWatch Logs via `CWL` client (not ECS API)
- Task definition loaded on demand (`FetchTaskDefinition`) — not cached at list level
- Breadcrumb status bar tracks the current drill-down path (cluster > service > task)
- Backspace/Esc navigates up one level in the hierarchy
- Each level has its own search, action menu, and detail overlay states
- Requires `enableExecuteCommand` on ECS service + SSM permissions on task role for ECS Exec
