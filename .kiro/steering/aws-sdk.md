# AWS SDK Conventions

- All AWS clients created in `internal/aws/session.go` NewClients factory
- 18 service clients in Clients struct: EC2, SSM, STS, ELBv2, ELB, ASG, CW, CWL, IAM, CF, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda
- Use AWS SDK paginators for list operations (FetchInstances, FetchSSMStatus)
- K8s integration: token via `aws eks get-token` (os/exec), direct HTTPS to EKS API server
- Profile parsing: handle both `[name]` (credentials) and `[profile name]` (config) formats
- "default" and InstanceRoleProfile omit `--profile` flag in SSM/ECS commands
- New AWS service: add client to Clients struct, create `<service>.go` in `internal/aws/`
