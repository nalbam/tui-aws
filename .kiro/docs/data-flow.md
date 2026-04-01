# Data Flow

## User Input → Rendering
1. User input → RootModel checks global keys (tab switch, profile, region, quit)
2. Non-global keys → delegated to active tab's Update()
3. Tab returns tea.Cmd for async AWS API calls
4. Response messages flow back through Update → re-render
5. SSM sessions: EC2 tab sends SSMExecRequest → RootModel runs tea.Exec (TUI suspended)
6. ECS Exec: ECS tab sends ECSExecRequest → RootModel runs tea.Exec
7. Tab switching: NavigateToTab message → RootModel switches active tab

## Caching
- Lazy loading: tabs fetch data on first activation
- 30-second TTL: fresh → use cache; stale → background reload showing cached data
- Cache invalidation: profile/region change → clear all; `R` key → clear active tab; SSM return → clear EC2 tab
- Memory only, no disk cache

## SSM Session Flow
```
EC2Tab → SSMExecRequest → RootModel → ssmExecCmd.Run()
  → stty sane + TCIFLUSH → SSMSessionDoneMsg
  → RootModel records history → Forward to EC2Tab → reload
```

## ECS Exec Flow
```
ECSTab → ECSExecRequest → RootModel → ecsExecCmd.Run()
  → stty sane + TCIFLUSH → ECSExecDoneMsg
  → Forward to ECSTab → resume container list
```

## K8s REST API (EKS Tab)
```
EKS Cluster → aws eks get-token (os/exec)
  → Bearer token → HTTPS to cluster endpoint
  → TLS via cluster CA certificate
```
