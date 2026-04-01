# Project Structure

```
tui-aws/
├── main.go                          Entry point, config migration, TUI launch
├── Makefile                         Build targets (build, build-all, test, clean)
├── scripts/setup.sh                 Cross-platform setup & install script
├── internal/
│   ├── config/config.go             Config load/save (~/.tui-aws/config.json)
│   ├── store/
│   │   ├── favorites.go             Favorites CRUD + persistence
│   │   └── history.go               Session history FIFO + persistence
│   ├── aws/
│   │   ├── session.go               SDK client factory (18 service clients)
│   │   ├── profile.go               AWS profile parsing
│   │   ├── ec2.go                   Instance model, FetchInstances
│   │   ├── ssm.go                   SSM command building
│   │   ├── vpc.go                   VPC, IGW, NAT, Peering, TGW, Endpoint, EIP
│   │   ├── subnet.go               Subnet, ENI
│   │   ├── network.go              RouteTable, Route entries
│   │   ├── security.go             SecurityGroup, NetworkACL rules
│   │   ├── reachability.go         VPC Reachability Analyzer
│   │   ├── elb.go                  ALB/NLB/CLB, listeners, targets
│   │   ├── asg.go                  Auto Scaling Groups
│   │   ├── ebs.go                  EBS volumes
│   │   ├── tgw.go                  Transit Gateways
│   │   ├── cloudwatch.go           CloudWatch alarms
│   │   ├── iam.go                  IAM users, groups, policies
│   │   ├── cloudfront.go           CloudFront distributions
│   │   ├── waf.go                  WAFv2 Web ACLs
│   │   ├── acm.go                  ACM certificates
│   │   ├── r53.go                  Route 53 hosted zones
│   │   ├── rds.go                  RDS DB instances
│   │   ├── s3.go                   S3 buckets
│   │   ├── ecs.go                  ECS clusters, services, tasks, exec
│   │   ├── eks.go                  EKS clusters, node groups
│   │   ├── k8s.go                  K8s REST API (no kubectl)
│   │   └── lambda.go               Lambda functions
│   └── ui/
│       ├── root.go                  RootModel, tab switching, global overlays
│       ├── tab.go                   Re-exports TabModel, SharedState, TabID
│       ├── placeholder.go           PlaceholderTab
│       ├── shared/                  TabModel interface, styles, table, overlay, selector
│       └── tab_*/                   22 tab packages (ec2, vpc, subnet, ...)
├── docs/
│   ├── architecture.md              System architecture
│   ├── decisions/                   ADRs
│   └── runbooks/                    Operational runbooks
└── .kiro/                           Kiro agent config, steering, docs, skills
```
