# Architecture

## System Overview
tui-aws is a single-binary Go CLI providing a 6-tab terminal UI for AWS infrastructure management. Built on Bubble Tea v2 (Elm architecture: Model-View-Update) with a tab-based architecture separating concerns into independent submodels.

## Tab Architecture

```
RootModel (tea.Model)
├── SharedState (profile, region, clients, cache, dimensions)
├── Tab Bar (1-6 switching, active highlight)
├── Global Overlays (profile/region selector)
└── TabModel[] (each implements Init/Update/View/ShortHelp)
    ├── [1] EC2Model      — instances, SSM, port forward, Network Path
    ├── [2] VPCModel      — VPCs, IGW, NAT, Peering, TGW, Endpoint, EIP
    ├── [3] SubnetModel   — subnets, ENIs
    ├── [4] RouteModel    — route tables, route entries
    ├── [5] SGModel       — security groups, NACLs (f toggles)
    └── [6] CheckModel    — connectivity checker, Reachability Analyzer
```

## Components

| Package | Path | Role |
|---------|------|------|
| **aws** | `internal/aws/` | All AWS SDK calls (EC2, SSM, STS). Single `ec2.Client` for VPC/Subnet/SG/NACL/RT/IGW/NAT/Peering/TGW/Endpoint/EIP/Reachability |
| **config** | `internal/config/` | User preferences (`~/.tui-aws/config.json`) |
| **store** | `internal/store/` | Favorites & session history persistence |
| **ui/root** | `internal/ui/root.go` | RootModel: tab switching, global keys, SSM exec, InterruptFilter |
| **ui/shared** | `internal/ui/shared/` | TabModel interface, SharedState, styles, table renderer, selector, overlay |
| **ui/tab_ec2** | `internal/ui/tab_ec2/` | EC2 tab: list, actions, search, filter, SSM, port forward, Network Path |
| **ui/tab_vpc** | `internal/ui/tab_vpc/` | VPC tab: list, details (lazy-loads sub-resources) |
| **ui/tab_subnet** | `internal/ui/tab_subnet/` | Subnet tab: list, ENI viewer |
| **ui/tab_routetable** | `internal/ui/tab_routetable/` | Route Table tab: list, route entries |
| **ui/tab_sg** | `internal/ui/tab_sg/` | SG/NACL tab: rules viewer, mode toggle |
| **ui/tab_troubleshoot** | `internal/ui/tab_troubleshoot/` | Connectivity checker engine + Reachability Analyzer |

## Data Flow

```
┌──────────┐      ┌──────────────┐      ┌────────────────┐
│  User    │─────▶│  RootModel   │─────▶│  Active Tab    │
│  Input   │      │  (global     │      │  .Update()     │
│  (keys)  │      │   keys)      │      └───────┬────────┘
└──────────┘      └──────┬───────┘              │
                         │                       │ tea.Cmd
                         │                       ▼
                    ┌────┴─────┐         ┌──────────────┐
                    │  View()  │◀────────│  AWS SDK      │
                    │  render  │         │  (async Cmd)  │
                    └──────────┘         └──────────────┘
```

1. User input → RootModel checks global keys (tab switch, profile, region, quit)
2. Non-global keys → delegated to active tab's Update()
3. Tab returns tea.Cmd for async AWS API calls
4. Response messages flow back through Update → re-render
5. SSM sessions: EC2 tab sends SSMExecRequest → RootModel runs tea.Exec (TUI suspended)
6. Tab switching: NavigateToTab message → RootModel switches active tab

## Caching Strategy

- **Lazy loading:** tabs fetch data on first activation, not at startup
- **30-second TTL:** fresh within TTL → use cache; stale → background reload showing cached data
- **Cache invalidation:** profile/region change → clear all; `R` key → clear active tab; SSM return → clear EC2 tab
- **Memory only:** no disk cache, SharedState.Cache map

## SSM Session Flow

```
EC2Tab → SSMExecRequest msg → RootModel intercepts
  → ssmExecCmd.Run() (aws ssm start-session)
  → stty sane + TCIFLUSH (terminal reset)
  → SSMSessionDoneMsg → RootModel records history
  → Forward to EC2Tab → reload instances
```

## Connectivity Checker (Tab 6)

Local 5-step validation:
1. Source SG Outbound
2. Source NACL Outbound
3. Source Route Table path
4. Destination NACL Inbound
5. Destination SG Inbound

Each step: ✓ pass / ✗ blocked (stops, shows suggestion). Optional AWS Reachability Analyzer API for confirmation.

## Infrastructure
- **Runtime:** Single binary, requires AWS CLI + Session Manager Plugin
- **Storage:** `~/.tui-aws/` (config.json, favorites.json, history.json)
- **Build:** Cross-compiled via Makefile for linux/darwin × amd64/arm64
- **Setup:** `scripts/setup.sh` auto-installs prerequisites on macOS/Linux
