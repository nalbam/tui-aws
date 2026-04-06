# tui-aws

![License](https://img.shields.io/badge/License-MIT-blue)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Version](https://img.shields.io/badge/Version-0.1.0-green)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)
[![English](https://img.shields.io/badge/lang-English-blue)](#english)
[![한국어](https://img.shields.io/badge/lang-한국어-red)](#한국어)

A terminal UI for exploring, managing, and troubleshooting entire AWS infrastructure — all from your terminal.

AWS 인프라 전체를 터미널에서 탐색, 관리, 트러블슈팅하는 TUI 도구.

---

# English

## Overview

**tui-aws** is a single-binary terminal UI that replaces the need to juggle multiple AWS Console tabs or remember complex CLI commands. It provides 22 integrated views covering EC2, VPC networking, containers, serverless, DNS, CDN, security, and more — with built-in SSM Session Manager, ECS Exec, and a local connectivity checker that validates SG + Route + NACL rules without calling AWS APIs.

[![Demo Video](https://img.youtube.com/vi/78gfU_Vfluw/maxresdefault.jpg)](https://youtu.be/78gfU_Vfluw)

> Click the image above to watch the demo video on YouTube

```
┌─ EC2 ASG EBS VPC Subnet Routes SG VPCE TGW ELB CF WAF ACM R53 RDS S3 ECS EKS Lambda CW IAM Check ─┐
│ ★ ● web-server-1         i-0abc1234   running  10.0.1.10   t3.medium  2a                        │
│   ● web-server-2         i-0def5678   running  10.0.1.11   t3.medium  2c                        │
│   ● db-primary           i-0ghi9012   running  10.0.2.20   r5.xlarge  2a                        │
│   ○ batch-worker         i-0jkl3456   stopped  10.0.3.30   c5.2xlarge 2b                        │
├──────────────────────────────────────────────────────────────────────────────────────────────────┤
│ ↑↓:Navigate  Enter:Actions  /:Search  f:Filter  p:Profile  r:Region  s:Sort  F:Fav  q:Quit     │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## Features

- **22 integrated tabs** — EC2, ASG, EBS, VPC, Subnet, Routes, SG/NACL, VPCE, TGW, ELB, CloudFront, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda, CloudWatch, IAM, Connectivity Check
- **SSM Session Manager** — Connect to instances without SSH keys or open security group rules. The TUI suspends, gives full terminal control, and resumes on exit.
- **ECS Exec** — Deep-dive into ECS clusters with Clusters > Services > Tasks > Containers > CloudWatch Logs > interactive shell via `aws ecs execute-command`
- **EKS K8s integration** — Browse Pods, Deployments, Services, and Nodes via direct K8s REST API calls (no kubectl dependency). Token from `aws eks get-token` with 14-min caching.
- **Local connectivity checker** — Validate SG + Route + NACL rules between any two EC2 instances in 5 steps without calling AWS APIs. Shows the exact blocking rule and a fix suggestion.
- **Network path visualization** — Trace VPC > Subnet > Route Table > Security Group > NACL in one scrollable overlay
- **Cross-resource navigation** — Drill down from EC2 to VPC, Subnet, Route Table, or Security Group with a single keystroke
- **Port forwarding** — Tunnel local ports to remote instances via SSM (RDS, web servers, debug ports)
- **Favorites and history** — Pin frequently accessed instances and track SSM session history

## Prerequisites

| Tool | Required | Purpose | Install |
|------|----------|---------|---------|
| **AWS CLI v2** | Yes | Runs `aws ssm start-session` and `aws ecs execute-command` | [Install guide](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) |
| **Session Manager Plugin** | Yes | Enables SSM session and ECS Exec connections | [Install guide](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) |
| **Go 1.25+** | Build only | Compiles the binary | [go.dev/dl](https://go.dev/dl/) |
| **AWS Credentials** | Yes | API access | `aws configure`, environment variables, or EC2 Instance Role |

Supported platforms: macOS (arm64, amd64), Linux (arm64, amd64).

## Installation

### One-line install and run

```bash
# Clone, install prerequisites, build, and launch
git clone https://github.com/whchoi98/tui-aws.git && cd tui-aws && ./scripts/setup.sh
```

The setup script checks your system, installs missing packages (with confirmation prompts), builds the binary, and launches tui-aws.

### Manual build (if prerequisites are already installed)

```bash
# Clone the repository
git clone https://github.com/whchoi98/tui-aws.git
cd tui-aws

# Build for current platform
make build

# Run
./tui-aws
```

### Cross-compile all platforms

```bash
# Build for linux/darwin x amd64/arm64
make build-all

# Output binaries in dist/
ls dist/
# tui-aws-linux-amd64
# tui-aws-linux-arm64
# tui-aws-darwin-arm64
# tui-aws-darwin-amd64
```

## Usage

### Basic usage

```bash
./tui-aws              # Launch TUI
./tui-aws --version    # Print version
```

### Switching profiles and regions

Press `p` to open the profile selector (named profiles from `~/.aws/credentials` and `~/.aws/config`, plus instance role). Press `r` to open the region selector. Changing either reloads the active tab's data.

### Search and filter

Press `/` in any tab to search by name, ID, or IP. Press `f` to open filters (EC2: state filter, SG: toggle SG/NACL mode). Press `Esc` to clear.

### SSM Session

1. Select an instance in the EC2 tab
2. Press `Enter` > select **SSM Session**
3. The TUI suspends and you get a full shell on the instance
4. Type `exit` or `Ctrl+D` to return to the TUI
5. Instance list refreshes automatically

### Port forwarding

1. Select an instance > `Enter` > **Port Forwarding**
2. Enter local port (default: 8080) and remote port (default: 80)
3. Press `Enter` to start the tunnel
4. Access the service at `localhost:<local-port>` from another terminal

### Connectivity check

1. Switch to the Check tab
2. Pick Source and Destination instances
3. Set Protocol (tcp/udp/all) and Port
4. Press `Enter` to run the 5-step local check
5. Optionally press `R` for AWS Reachability Analyzer (may incur costs)

```
  Connectivity: web-server → db-primary  TCP/443
  ══════════════════════════════════════════════

  ✓ Source SG Outbound     sg-0abc: TCP 443 → 0.0.0.0/0 ALLOW
  ✓ Source NACL Outbound   acl-xxx: Rule 100 All ALLOW
  ✓ Source Route           rtb-xxx: 10.2.0.0/16 → tgw-xxx (active)
  ✗ Dest SG Inbound        sg-0def: TCP 443 ← 10.1.0.0/16 NOT FOUND

  Result: ✗ BLOCKED at Destination SG Inbound
  Suggestion: Add inbound rule TCP 443 from 10.1.88.66/32
```

### Key bindings

**Global keys (all tabs):**

| Key | Action |
|-----|--------|
| `]` / `[` | Next / previous tab |
| `Tab` / `Shift+Tab` | Next / previous tab |
| `p` | Select AWS profile |
| `r` | Select AWS region |
| `R` | Refresh current tab data |
| `q` / `Ctrl+C` | Quit |

**Table keys (all tabs):**

| Key | Action |
|-----|--------|
| `Up` `Down` / `j` `k` | Move cursor |
| `Enter` | Open action menu |
| `/` | Start search |
| `f` | Open filter / toggle mode |
| `s` / `S` | Cycle sort column / reverse direction |
| `F` | Toggle favorite (EC2 tab) |
| `Esc` | Close overlay, cancel search |

### Tabs reference

| Tab | Description |
|-----|-------------|
| **EC2** | Instances with SSM, port forward, Network Path, favorites |
| **ASG** | Auto Scaling Groups, scaling policies, instances |
| **EBS** | Volumes with encryption status (color-coded) |
| **VPC** | VPCs with details (IGW, NAT, Peering, TGW, Endpoint, EIP) |
| **Subnet** | Subnets with ENI viewer |
| **Routes** | Route tables with route entries |
| **SG** | Security Groups / NACLs (toggle with `f`) |
| **VPCE** | VPC Endpoints (Gateway/Interface) |
| **TGW** | Transit Gateways with attachments and routes |
| **ELB** | ALB/NLB/CLB with interactive target group detail |
| **CF** | CloudFront distributions |
| **WAF** | WAFv2 Web ACLs with rules and associated resources |
| **ACM** | Certificates with status, expiry, SANs |
| **R53** | Route 53 hosted zones with records (on demand) |
| **RDS** | DB instances with engine, class, endpoint |
| **S3** | Buckets with versioning, encryption, public access |
| **ECS** | Clusters > Services > Tasks > Containers > Logs > ECS Exec |
| **EKS** | Clusters > Namespaces > Pods/Deployments/Services, Nodes, Pod Logs |
| **Lambda** | Functions with runtime, memory, VPC config, layers |
| **CW** | CloudWatch alarms with state, metric, threshold |
| **IAM** | Users with groups, policies, last used |
| **Check** | Connectivity checker + Reachability Analyzer |

## Configuration

All config files are stored in `~/.tui-aws/` (created automatically on first run).

| File | Purpose |
|------|---------|
| `config.json` | Default profile, region, table display settings |
| `favorites.json` | Favorited instances, keyed by instance ID + profile + region |
| `history.json` | SSM session history, FIFO with max 100 entries |

### config.json

```json
{
  "default_profile": "default",
  "default_region": "ap-northeast-2",
  "refresh_interval_seconds": 0,
  "table": {
    "visible_columns": ["name", "id", "state", "private_ip", "type", "az"],
    "sort_by": "name",
    "sort_order": "asc"
  }
}
```

| Field | Default | Description |
|-------|---------|-------------|
| `default_profile` | `"default"` | AWS profile to use on startup |
| `default_region` | `"us-east-1"` | AWS region to use on startup |
| `refresh_interval_seconds` | `0` | Auto-refresh interval (0 = manual only) |
| `table.sort_by` | `"name"` | Default sort column |
| `table.sort_order` | `"asc"` | Default sort direction |

## Project Structure

```
tui-aws/
├── main.go                          # Entry point, config migration, TUI launch
├── Makefile                         # Build targets (build, build-all, install, test, clean)
├── scripts/
│   └── setup.sh                     # Cross-platform setup and install script
├── internal/
│   ├── aws/                         # AWS SDK integration (18 service clients + K8s REST)
│   │   ├── ec2.go                   # Instance model, FetchInstances
│   │   ├── vpc.go                   # VPC, IGW, NAT, Peering, TGW, Endpoint, EIP
│   │   ├── subnet.go               # Subnet, ENI
│   │   ├── network.go              # Route Table, Route entries
│   │   ├── security.go             # Security Group rules, Network ACL rules
│   │   ├── reachability.go         # VPC Reachability Analyzer
│   │   ├── profile.go              # AWS profile parsing
│   │   ├── session.go              # SDK client factory (18 clients)
│   │   ├── ssm.go                  # SSM command building
│   │   ├── elb.go                  # ALB/NLB/CLB, listeners, targets
│   │   ├── asg.go                  # Auto Scaling Groups
│   │   ├── ebs.go                  # EBS volumes
│   │   ├── tgw.go                  # Transit Gateways
│   │   ├── cloudwatch.go           # CloudWatch alarms
│   │   ├── iam.go                  # IAM users, groups, policies
│   │   ├── cloudfront.go           # CloudFront distributions
│   │   ├── waf.go                  # WAFv2 Web ACLs
│   │   ├── acm.go                  # ACM certificates
│   │   ├── r53.go                  # Route 53 hosted zones, records
│   │   ├── rds.go                  # RDS DB instances
│   │   ├── s3.go                   # S3 buckets
│   │   ├── ecs.go                  # ECS clusters, services, tasks, exec
│   │   ├── eks.go                  # EKS clusters, node groups
│   │   ├── k8s.go                  # K8s REST API (no kubectl)
│   │   └── lambda.go               # Lambda functions
│   ├── config/
│   │   └── config.go               # Load/save user config (~/.tui-aws/config.json)
│   ├── store/
│   │   ├── favorites.go            # Favorites CRUD + persistence
│   │   └── history.go              # Session history FIFO + persistence
│   └── ui/
│       ├── root.go                  # RootModel, tab switching, SSM/ECS exec
│       ├── tab.go                   # Re-exports TabModel, SharedState, TabID
│       ├── shared/                  # TabModel interface, styles, table, overlay
│       ├── tab_ec2/                 # EC2 tab (6 files)
│       ├── tab_asg/ ... tab_iam/   # 20 additional tab packages (3 files each)
│       └── tab_troubleshoot/        # Connectivity checker (4 files)
├── docs/
│   ├── architecture.md              # System architecture
│   ├── decisions/                   # Architecture Decision Records
│   ├── runbooks/                    # Operational runbooks
│   └── onboarding.md               # Developer onboarding guide
└── .claude/                         # Claude Code hooks, skills, commands, agents
```

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Static analysis
go vet ./...

# Test a specific package
go test ./internal/ui/tab_troubleshoot/ -v
```

## IAM Permissions

### Minimum (EC2 + SSM only)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:DescribeInstances",
      "ec2:DescribeVpcs",
      "ec2:DescribeSubnets",
      "ssm:StartSession",
      "ssm:TerminateSession",
      "ssm:DescribeInstanceInformation",
      "sts:GetCallerIdentity"
    ],
    "Resource": "*"
  }]
}
```

### Full (all 22 tabs)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:Describe*",
      "ec2:SearchTransitGatewayRoutes",
      "ssm:StartSession",
      "ssm:TerminateSession",
      "ssm:DescribeInstanceInformation",
      "sts:GetCallerIdentity",
      "autoscaling:DescribeAutoScalingGroups",
      "elasticloadbalancing:Describe*",
      "cloudfront:ListDistributions",
      "wafv2:ListWebACLs",
      "wafv2:GetWebACL",
      "wafv2:ListResourcesForWebACL",
      "acm:ListCertificates",
      "acm:DescribeCertificate",
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
      "rds:DescribeDBInstances",
      "s3:ListAllMyBuckets",
      "s3:GetBucket*",
      "ecs:List*",
      "ecs:Describe*",
      "ecs:ExecuteCommand",
      "eks:ListClusters",
      "eks:DescribeCluster",
      "eks:ListNodegroups",
      "eks:DescribeNodegroup",
      "lambda:ListFunctions",
      "cloudwatch:DescribeAlarms",
      "logs:GetLogEvents",
      "iam:ListUsers",
      "iam:ListGroupsForUser",
      "iam:ListAttachedUserPolicies"
    ],
    "Resource": "*"
  }]
}
```

### Reachability Analyzer (optional, may incur costs)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:CreateNetworkInsightsPath",
      "ec2:DeleteNetworkInsightsPath",
      "ec2:StartNetworkInsightsAnalysis",
      "ec2:DescribeNetworkInsightsAnalyses"
    ],
    "Resource": "*"
  }]
}
```

> **Note:** If a tab shows "AccessDenied", only that tab is affected. Other tabs continue working.

## Troubleshooting

### "exit status 255" when connecting via SSM

| Cause | Solution |
|-------|----------|
| Invalid AWS credentials | Check `~/.aws/credentials` for syntax errors |
| Missing SSM Agent | Verify the instance has SSM Agent installed and running |
| Missing IAM role | Attach `AmazonSSMManagedInstanceCore` policy to the instance role |
| VPC endpoint missing | For private subnets without NAT, create SSM VPC endpoints (`ssm`, `ssmmessages`, `ec2messages`) |
| Wrong profile/region | Press `p`/`r` to switch in tui-aws |

### "AccessDenied" on a tab

The current IAM identity lacks the required permissions. See [IAM Permissions](#iam-permissions) for the full policy. Only the affected tab shows the error.

### Garbled text or broken columns

Ensure your terminal supports UTF-8, 256-color or TrueColor, and a monospace font with Unicode support (JetBrains Mono, Fira Code, Menlo). For SSH: `export TERM=xterm-256color`.

### TUI does not return after SSM session

tui-aws includes terminal reset (`stty sane` + stdin flush) after SSM sessions. If issues persist, run `reset` or `stty sane` manually.

## Contributing

1. Fork the repository
2. Create a feature branch
   ```bash
   git checkout -b feat/your-feature
   ```
3. Commit with Conventional Commits format
   ```bash
   git commit -m "feat: add support for new AWS service"
   git commit -m "fix: resolve nil pointer in ECS tab"
   ```
4. Push to your fork
   ```bash
   git push origin feat/your-feature
   ```
5. Open a Pull Request against `main`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

- **Maintainer:** whchoi98 — [GitHub](https://github.com/whchoi98) — [whchoi98@gmail.com](mailto:whchoi98@gmail.com)
- **Issues:** [github.com/whchoi98/tui-aws/issues](https://github.com/whchoi98/tui-aws/issues)

---

# 한국어

## 개요

**tui-aws**는 여러 AWS 콘솔 탭을 오가거나 복잡한 CLI 명령을 기억할 필요 없이, 터미널 하나에서 AWS 인프라를 관리할 수 있는 단일 바이너리 도구입니다. EC2, VPC 네트워킹, 컨테이너, 서버리스, DNS, CDN, 보안 등 22개 통합 뷰를 제공하며, SSM Session Manager, ECS Exec, 그리고 AWS API 호출 없이 SG + Route + NACL 규칙을 검증하는 로컬 연결성 검사기를 내장하고 있습니다.

[![데모 영상](https://img.youtube.com/vi/78gfU_Vfluw/maxresdefault.jpg)](https://youtu.be/78gfU_Vfluw)

> 위 이미지를 클릭하면 YouTube 데모 영상으로 이동합니다

```
┌─ EC2 ASG EBS VPC Subnet Routes SG VPCE TGW ELB CF WAF ACM R53 RDS S3 ECS EKS Lambda CW IAM Check ─┐
│ ★ ● web-server-1         i-0abc1234   running  10.0.1.10   t3.medium  2a                        │
│   ● web-server-2         i-0def5678   running  10.0.1.11   t3.medium  2c                        │
│   ● db-primary           i-0ghi9012   running  10.0.2.20   r5.xlarge  2a                        │
│   ○ batch-worker         i-0jkl3456   stopped  10.0.3.30   c5.2xlarge 2b                        │
├──────────────────────────────────────────────────────────────────────────────────────────────────┤
│ ↑↓:Navigate  Enter:Actions  /:Search  f:Filter  p:Profile  r:Region  s:Sort  F:Fav  q:Quit     │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## 주요 기능

- **22개 통합 탭** — EC2, ASG, EBS, VPC, Subnet, Routes, SG/NACL, VPCE, TGW, ELB, CloudFront, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda, CloudWatch, IAM, 연결성 검사
- **SSM Session Manager** — SSH 키나 보안 그룹 인바운드 규칙 없이 인스턴스에 접속합니다. TUI가 일시 중지되고 터미널 제어가 넘어가며, 종료 시 자동 복귀합니다.
- **ECS Exec** — ECS 클러스터를 깊이 탐색합니다. Clusters > Services > Tasks > Containers > CloudWatch Logs > `aws ecs execute-command`를 통한 대화형 셸을 제공합니다.
- **EKS K8s 통합** — K8s REST API를 직접 호출하여 Pods, Deployments, Services, Nodes를 탐색합니다 (kubectl 불필요). `aws eks get-token`으로 토큰을 생성하고 14분간 캐싱합니다.
- **로컬 연결성 검사기** — AWS API 호출 없이 두 EC2 인스턴스 간 SG + Route + NACL 규칙을 5단계로 검증합니다. 정확한 차단 규칙과 수정 제안을 표시합니다.
- **네트워크 경로 시각화** — VPC > Subnet > Route Table > Security Group > NACL을 하나의 스크롤 가능한 오버레이에서 확인합니다.
- **크로스 리소스 탐색** — EC2에서 VPC, Subnet, Route Table, Security Group으로 한 번의 키 입력으로 이동합니다.
- **포트 포워딩** — SSM을 통해 로컬 포트를 원격 인스턴스로 터널링합니다 (RDS, 웹 서버, 디버그 포트).
- **즐겨찾기와 이력** — 자주 접근하는 인스턴스를 고정하고 SSM 세션 이력을 추적합니다.

## 사전 요구 사항

| 도구 | 필수 | 용도 | 설치 |
|------|------|------|------|
| **AWS CLI v2** | 예 | `aws ssm start-session` 및 `aws ecs execute-command` 실행 | [설치 가이드](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) |
| **Session Manager Plugin** | 예 | SSM 세션 및 ECS Exec 연결 | [설치 가이드](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) |
| **Go 1.25+** | 빌드 시 | 바이너리 컴파일 | [go.dev/dl](https://go.dev/dl/) |
| **AWS 자격 증명** | 예 | API 접근 | `aws configure`, 환경변수, 또는 EC2 Instance Role |

지원 플랫폼: macOS (arm64, amd64), Linux (arm64, amd64).

## 설치 방법

### 한 줄로 설치 및 실행

```bash
# 클론, 사전 요구 사항 설치, 빌드, 실행
git clone https://github.com/whchoi98/tui-aws.git && cd tui-aws && ./scripts/setup.sh
```

설치 스크립트가 시스템을 점검하고, 부족한 패키지를 설치하고(확인 프롬프트 표시), 바이너리를 빌드한 후 tui-aws를 실행합니다.

### 수동 빌드 (사전 요구 사항이 이미 설치된 경우)

```bash
# 저장소 클론
git clone https://github.com/whchoi98/tui-aws.git
cd tui-aws

# 현재 플랫폼 빌드
make build

# 실행
./tui-aws
```

### 크로스 컴파일

```bash
# linux/darwin x amd64/arm64 빌드
make build-all

# 빌드 결과물 확인
ls dist/
# tui-aws-linux-amd64
# tui-aws-linux-arm64
# tui-aws-darwin-arm64
# tui-aws-darwin-amd64
```

## 사용법

### 기본 사용

```bash
./tui-aws              # TUI 실행
./tui-aws --version    # 버전 출력
```

### 프로파일 및 리전 변경

`p`를 눌러 프로파일 선택기를 엽니다 (`~/.aws/credentials`와 `~/.aws/config`의 명명된 프로파일 + 인스턴스 역할). `r`을 눌러 리전 선택기를 엽니다. 변경 시 현재 탭의 데이터가 리로드됩니다.

### 검색 및 필터링

아무 탭에서 `/`를 눌러 이름, ID, IP로 검색합니다. `f`를 눌러 필터를 엽니다 (EC2: 상태 필터, SG: SG/NACL 모드 전환). `Esc`로 해제합니다.

### SSM 세션

1. EC2 탭에서 인스턴스를 선택합니다
2. `Enter` > **SSM Session**을 선택합니다
3. TUI가 일시 중지되고 인스턴스에서 전체 셸을 사용합니다
4. `exit` 또는 `Ctrl+D`로 TUI에 복귀합니다
5. 인스턴스 목록이 자동으로 새로고침됩니다

### 포트 포워딩

1. 인스턴스 선택 > `Enter` > **Port Forwarding**
2. 로컬 포트(기본: 8080)와 리모트 포트(기본: 80)를 입력합니다
3. `Enter`로 터널을 시작합니다
4. 다른 터미널에서 `localhost:<로컬포트>`로 서비스에 접근합니다

### 연결성 검사

1. Check 탭으로 이동합니다
2. Source와 Destination 인스턴스를 선택합니다
3. Protocol (tcp/udp/all)과 Port를 설정합니다
4. `Enter`로 5단계 로컬 검사를 실행합니다
5. 선택적으로 `R`을 눌러 AWS Reachability Analyzer를 실행합니다 (비용 발생 가능)

```
  Connectivity: web-server → db-primary  TCP/443
  ══════════════════════════════════════════════

  ✓ Source SG Outbound     sg-0abc: TCP 443 → 0.0.0.0/0 ALLOW
  ✓ Source NACL Outbound   acl-xxx: Rule 100 All ALLOW
  ✓ Source Route           rtb-xxx: 10.2.0.0/16 → tgw-xxx (active)
  ✗ Dest SG Inbound        sg-0def: TCP 443 ← 10.1.0.0/16 NOT FOUND

  Result: ✗ BLOCKED at Destination SG Inbound
  Suggestion: Add inbound rule TCP 443 from 10.1.88.66/32
```

### 키 바인딩

**전역 키 (모든 탭):**

| 키 | 동작 |
|----|------|
| `]` / `[` | 다음 / 이전 탭 |
| `Tab` / `Shift+Tab` | 다음 / 이전 탭 |
| `p` | AWS 프로파일 선택 |
| `r` | AWS 리전 선택 |
| `R` | 현재 탭 데이터 새로고침 |
| `q` / `Ctrl+C` | 종료 |

**테이블 키 (모든 탭):**

| 키 | 동작 |
|----|------|
| `Up` `Down` / `j` `k` | 커서 이동 |
| `Enter` | 액션 메뉴 열기 |
| `/` | 검색 시작 |
| `f` | 필터 열기 / 모드 전환 |
| `s` / `S` | 정렬 컬럼 순환 / 방향 반전 |
| `F` | 즐겨찾기 토글 (EC2 탭) |
| `Esc` | 오버레이 닫기, 검색 취소 |

### 탭 참조

| 탭 | 설명 |
|----|------|
| **EC2** | 인스턴스 — SSM, 포트 포워딩, Network Path, 즐겨찾기 |
| **ASG** | Auto Scaling Groups — 스케일링 정책, 인스턴스 |
| **EBS** | 볼륨 — 암호화 상태 (색상 표시) |
| **VPC** | VPC 상세 (IGW, NAT, Peering, TGW, Endpoint, EIP) |
| **Subnet** | 서브넷 — ENI 뷰어 |
| **Routes** | 라우트 테이블 — 경로 엔트리 |
| **SG** | Security Groups / NACLs (`f`로 전환) |
| **VPCE** | VPC Endpoints (Gateway/Interface) |
| **TGW** | Transit Gateways — 어태치먼트, 라우트 |
| **ELB** | ALB/NLB/CLB — 대화형 타겟 그룹 상세 |
| **CF** | CloudFront 배포 |
| **WAF** | WAFv2 Web ACLs — 규칙, 연결 리소스 |
| **ACM** | 인증서 — 상태, 만료일, SANs |
| **R53** | Route 53 호스팅 존 — 레코드 (온디맨드) |
| **RDS** | DB 인스턴스 — 엔진, 클래스, 엔드포인트 |
| **S3** | 버킷 — 버전관리, 암호화, 퍼블릭 접근 |
| **ECS** | Clusters > Services > Tasks > Containers > Logs > ECS Exec |
| **EKS** | Clusters > Namespaces > Pods/Deployments/Services, Nodes, Pod Logs |
| **Lambda** | 함수 — 런타임, 메모리, VPC 설정, 레이어 |
| **CW** | CloudWatch 알람 — 상태, 메트릭, 임계값 |
| **IAM** | 사용자 — 그룹, 정책, 마지막 사용 |
| **Check** | 연결성 검사 + Reachability Analyzer |

## 환경 설정

모든 설정 파일은 `~/.tui-aws/`에 저장됩니다 (첫 실행 시 자동 생성).

| 파일 | 용도 |
|------|------|
| `config.json` | 기본 프로파일, 리전, 테이블 표시 설정 |
| `favorites.json` | 즐겨찾기 인스턴스 (instance ID + profile + region으로 키) |
| `history.json` | SSM 세션 이력 (최대 100개 FIFO) |

### config.json

```json
{
  "default_profile": "default",
  "default_region": "ap-northeast-2",
  "refresh_interval_seconds": 0,
  "table": {
    "visible_columns": ["name", "id", "state", "private_ip", "type", "az"],
    "sort_by": "name",
    "sort_order": "asc"
  }
}
```

| 필드 | 기본값 | 설명 |
|------|--------|------|
| `default_profile` | `"default"` | 시작 시 사용할 AWS 프로파일 |
| `default_region` | `"us-east-1"` | 시작 시 사용할 AWS 리전 |
| `refresh_interval_seconds` | `0` | 자동 새로고침 간격 (0 = 수동만) |
| `table.sort_by` | `"name"` | 기본 정렬 컬럼 |
| `table.sort_order` | `"asc"` | 기본 정렬 방향 |

## 프로젝트 구조

```
tui-aws/
├── main.go                          # 진입점, 설정 마이그레이션, TUI 실행
├── Makefile                         # 빌드 타겟 (build, build-all, install, test, clean)
├── scripts/
│   └── setup.sh                     # 크로스 플랫폼 설치 스크립트
├── internal/
│   ├── aws/                         # AWS SDK 통합 (18개 서비스 클라이언트 + K8s REST)
│   │   ├── ec2.go                   # 인스턴스 모델, FetchInstances
│   │   ├── vpc.go                   # VPC, IGW, NAT, Peering, TGW, Endpoint, EIP
│   │   ├── subnet.go               # Subnet, ENI
│   │   ├── network.go              # Route Table, Route 엔트리
│   │   ├── security.go             # Security Group 규칙, Network ACL 규칙
│   │   ├── reachability.go         # VPC Reachability Analyzer
│   │   ├── profile.go              # AWS 프로파일 파싱
│   │   ├── session.go              # SDK 클라이언트 팩토리 (18개 클라이언트)
│   │   ├── ssm.go                  # SSM 명령 빌더
│   │   ├── elb.go                  # ALB/NLB/CLB, 리스너, 타겟
│   │   ├── asg.go                  # Auto Scaling Groups
│   │   ├── ebs.go                  # EBS 볼륨
│   │   ├── tgw.go                  # Transit Gateways
│   │   ├── cloudwatch.go           # CloudWatch 알람
│   │   ├── iam.go                  # IAM 사용자, 그룹, 정책
│   │   ├── cloudfront.go           # CloudFront 배포
│   │   ├── waf.go                  # WAFv2 Web ACLs
│   │   ├── acm.go                  # ACM 인증서
│   │   ├── r53.go                  # Route 53 호스팅 존, 레코드
│   │   ├── rds.go                  # RDS DB 인스턴스
│   │   ├── s3.go                   # S3 버킷
│   │   ├── ecs.go                  # ECS 클러스터, 서비스, 태스크, exec
│   │   ├── eks.go                  # EKS 클러스터, 노드 그룹
│   │   ├── k8s.go                  # K8s REST API (kubectl 불필요)
│   │   └── lambda.go               # Lambda 함수
│   ├── config/
│   │   └── config.go               # 사용자 설정 로드/저장 (~/.tui-aws/config.json)
│   ├── store/
│   │   ├── favorites.go            # 즐겨찾기 CRUD + 영속화
│   │   └── history.go              # 세션 이력 FIFO + 영속화
│   └── ui/
│       ├── root.go                  # RootModel, 탭 전환, SSM/ECS exec
│       ├── tab.go                   # TabModel, SharedState, TabID 재수출
│       ├── shared/                  # TabModel 인터페이스, 스타일, 테이블, 오버레이
│       ├── tab_ec2/                 # EC2 탭 (6개 파일)
│       ├── tab_asg/ ... tab_iam/   # 20개 추가 탭 패키지 (각 3개 파일)
│       └── tab_troubleshoot/        # 연결성 검사기 (4개 파일)
├── docs/
│   ├── architecture.md              # 시스템 아키텍처
│   ├── decisions/                   # 아키텍처 의사결정 기록 (ADR)
│   ├── runbooks/                    # 운영 런북
│   └── onboarding.md               # 개발자 온보딩 가이드
└── .claude/                         # Claude Code 훅, 스킬, 명령, 에이전트
```

## 테스트

```bash
# 전체 테스트 실행
make test

# 상세 출력으로 테스트
go test ./... -v

# 정적 분석
go vet ./...

# 특정 패키지 테스트
go test ./internal/ui/tab_troubleshoot/ -v
```

## IAM 권한 설정

### 최소 (EC2 + SSM만)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:DescribeInstances",
      "ec2:DescribeVpcs",
      "ec2:DescribeSubnets",
      "ssm:StartSession",
      "ssm:TerminateSession",
      "ssm:DescribeInstanceInformation",
      "sts:GetCallerIdentity"
    ],
    "Resource": "*"
  }]
}
```

### 전체 (22개 탭)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:Describe*",
      "ec2:SearchTransitGatewayRoutes",
      "ssm:StartSession",
      "ssm:TerminateSession",
      "ssm:DescribeInstanceInformation",
      "sts:GetCallerIdentity",
      "autoscaling:DescribeAutoScalingGroups",
      "elasticloadbalancing:Describe*",
      "cloudfront:ListDistributions",
      "wafv2:ListWebACLs",
      "wafv2:GetWebACL",
      "wafv2:ListResourcesForWebACL",
      "acm:ListCertificates",
      "acm:DescribeCertificate",
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
      "rds:DescribeDBInstances",
      "s3:ListAllMyBuckets",
      "s3:GetBucket*",
      "ecs:List*",
      "ecs:Describe*",
      "ecs:ExecuteCommand",
      "eks:ListClusters",
      "eks:DescribeCluster",
      "eks:ListNodegroups",
      "eks:DescribeNodegroup",
      "lambda:ListFunctions",
      "cloudwatch:DescribeAlarms",
      "logs:GetLogEvents",
      "iam:ListUsers",
      "iam:ListGroupsForUser",
      "iam:ListAttachedUserPolicies"
    ],
    "Resource": "*"
  }]
}
```

### Reachability Analyzer (선택적, 비용 발생 가능)

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:CreateNetworkInsightsPath",
      "ec2:DeleteNetworkInsightsPath",
      "ec2:StartNetworkInsightsAnalysis",
      "ec2:DescribeNetworkInsightsAnalyses"
    ],
    "Resource": "*"
  }]
}
```

> **참고:** 탭에서 "AccessDenied"가 표시되면 해당 탭만 영향을 받습니다. 다른 탭은 정상 동작합니다.

## 문제 해결

### SSM 접속 시 "exit status 255"

| 원인 | 해결 방법 |
|------|----------|
| 잘못된 AWS 자격 증명 | `~/.aws/credentials`에서 구문 오류를 확인합니다 |
| 인스턴스에 SSM Agent 없음 | SSM Agent가 설치되고 실행 중인지 확인합니다 |
| IAM 역할 없음 | `AmazonSSMManagedInstanceCore` 정책을 인스턴스 역할에 연결합니다 |
| VPC 엔드포인트 없음 | NAT 없는 프라이빗 서브넷은 SSM VPC 엔드포인트를 생성합니다 (`ssm`, `ssmmessages`, `ec2messages`) |
| 잘못된 프로파일/리전 | tui-aws에서 `p`/`r`로 변경합니다 |

### 탭에서 "AccessDenied"

현재 IAM 자격 증명에 필요한 권한이 없습니다. [IAM 권한 설정](#iam-권한-설정)에서 전체 정책을 확인합니다. 해당 탭만 에러를 표시합니다.

### 텍스트 깨짐 또는 컬럼 정렬 오류

터미널이 UTF-8, 256색 또는 TrueColor, Unicode를 지원하는 고정폭 글꼴(JetBrains Mono, Fira Code, Menlo)을 지원하는지 확인합니다. SSH 사용 시: `export TERM=xterm-256color`.

### SSM 세션 후 TUI가 복귀하지 않음

tui-aws는 SSM 세션 후 터미널 리셋(`stty sane` + stdin flush)을 수행합니다. 문제가 지속되면 `reset` 또는 `stty sane`을 수동 실행합니다.

## 기여 방법

1. 저장소를 Fork합니다
2. 기능 브랜치를 생성합니다
   ```bash
   git checkout -b feat/your-feature
   ```
3. Conventional Commits 형식으로 커밋합니다
   ```bash
   git commit -m "feat: 새로운 AWS 서비스 지원 추가"
   git commit -m "fix: ECS 탭에서 nil 포인터 수정"
   ```
4. Fork한 저장소에 Push합니다
   ```bash
   git push origin feat/your-feature
   ```
5. `main` 브랜치에 대한 Pull Request를 생성합니다

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조합니다.

## 연락처

- **메인테이너:** whchoi98 — [GitHub](https://github.com/whchoi98) — [whchoi98@gmail.com](mailto:whchoi98@gmail.com)
- **이슈:** [github.com/whchoi98/tui-aws/issues](https://github.com/whchoi98/tui-aws/issues)
