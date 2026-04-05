# EKS Tab
## Role
EKS K8s integration — Clusters → Namespaces → Pods/Deployments/Services, Nodes, Pod Logs.
## Key Files
- `model.go` — 17-state viewState for K8s drill-down, breadcrumb status bar
- `table.go` — Per-level columns (cluster, namespace, pod, deployment, service, node)
- `detail.go` — Pod/deployment/service/node detail + pod log viewer
## Drill-Down Hierarchy
```
ClusterList → ClusterAction → NamespaceList → ResourceMenu → PodList / DeployList / ServiceList
                  ↓                                               ↓           ↓            ↓
             ClusterDetail                                   PodDetail  DeployDetail  ServiceDetail
              NodeGroupList → NodeList → NodeDetail                ↓
                                                               PodLogs
```
## Rules
- K8s REST API via `net/http` — NO client-go or kubectl dependency
- Token obtained via `aws eks get-token` (os/exec) with 14-minute caching in `internal/aws/k8s.go`
- TLS verified with cluster CA cert (base64 decoded from EKS DescribeCluster response)
- ResourceMenu presents Pods/Deployments/Services choice after selecting a namespace
- Pod logs: `/api/v1/namespaces/{ns}/pods/{name}/log?tailLines=50`
- Node groups fetched from AWS EKS API; K8s nodes fetched from K8s API (different sources)
- Breadcrumb status bar tracks: cluster > namespace > resource type
- Backspace/Esc navigates up one level
- Each level has its own async message types (clustersLoadedMsg, podsLoadedMsg, etc.)
