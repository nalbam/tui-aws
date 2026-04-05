# ADR-001: K8s REST API without client-go

## Status
Accepted

## Context
The EKS tab needs to display K8s resources (pods, deployments, services, nodes) and stream pod logs. The standard Go approach is to use `client-go`, but this adds a massive dependency tree (~50+ packages) and significantly increases binary size for what amounts to a handful of REST API calls.

## Decision
Use raw `net/http` calls to the K8s API server with bearer token authentication. Tokens are obtained via `aws eks get-token` (os/exec) and cached for 14 minutes. TLS is verified using the cluster CA certificate (base64-decoded from the EKS DescribeCluster response).

Implemented in `internal/aws/k8s.go`.

## Consequences

### Positive
- Single binary stays small (~25MB vs ~40MB+ with client-go)
- No kubectl dependency required at runtime
- Token management is simple and auditable
- Zero K8s dependency in go.mod

### Negative
- Must manually handle K8s API pagination and error responses
- No automatic client-side caching or informer pattern
- New K8s resource types require manual JSON struct definitions

### Risks
- K8s API version changes require manual struct updates (mitigated: using stable v1 APIs only)
