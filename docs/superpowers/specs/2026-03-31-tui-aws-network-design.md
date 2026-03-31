# tui-aws Network Extension Design Spec

AWS EC2 인스턴스 관리 도구(tui-ssm)를 **tui-aws**로 확장하여, VPC 네트워크 인프라 전체를 TUI에서 탐색하고 연결 문제를 진단할 수 있도록 한다.

## 이름 변경

- 모듈: `tui-ssm` → `tui-aws`
- 바이너리: `tui-aws`
- 설정 디렉토리: `~/.tui-ssm/` → `~/.tui-aws/` (기존 파일 자동 마이그레이션)

## 탭 구조

```
[1:EC2] [2:VPC] [3:Subnets] [4:Routes] [5:SG] [6:Check]
```

숫자키(`1`-`6`) 또는 Tab/Shift+Tab으로 전환. 활성 탭 하이라이트.

| 탭 | 키 | 핵심 기능 |
|----|-----|----------|
| EC2 | `1` | 기존 기능 + Network Path 오버레이 + Go to VPC/Subnet 드릴다운 |
| VPC | `2` | VPC 목록 + 상세(IGW, NAT, Peering, TGW, Endpoint, EIP) |
| Subnets | `3` | Subnet 목록(Name, CIDR, AZ, 가용IP, Public/Private) + ENI |
| Routes | `4` | Route Table 목록 + Route Entries(Dest→Target→State) |
| SG | `5` | Security Group 규칙 + Network ACL 규칙 (f키로 전환) |
| Check | `6` | 소스→대상 연결 검증 (SG+Route+NACL) + Reachability Analyzer 옵션 |

## 아키텍처: 탭별 서브모델

```
RootModel
├── 공유 상태: SharedState (profile, region, clients, width, height, cache)
├── 탭 상태: activeTab (0-5)
├── StatusBar: 탭 바 + profile/region
├── HelpBar: 전역키 + 탭별 키
└── TabModel[]
    ├── EC2Model      (기존 model.go 리팩토링)
    ├── VPCModel
    ├── SubnetModel
    ├── RouteModel
    ├── SGModel
    └── TroubleshootModel
```

### 서브모델 인터페이스

```go
type TabModel interface {
    Init(shared *SharedState) tea.Cmd
    Update(msg tea.Msg, shared *SharedState) (TabModel, tea.Cmd)
    View(shared *SharedState) string
    ShortHelp() string
}

type SharedState struct {
    Profile   string
    Region    string
    Profiles  []string
    Clients   *aws.Clients
    Width     int
    Height    int
    Favorites *store.Favorites
    History   *store.History
    Cache     map[string]CachedData
}

type CachedData struct {
    Data      any
    FetchedAt time.Time
}
```

### 프로젝트 구조

```
tui-aws/
├── main.go
├── internal/
│   ├── aws/
│   │   ├── ec2.go           # 기존
│   │   ├── profile.go       # 기존
│   │   ├── session.go       # 기존
│   │   ├── ssm.go           # 기존
│   │   ├── vpc.go           # VPC, IGW, NAT, Peering, TGW, Endpoint, EIP
│   │   ├── subnet.go        # Subnet, ENI
│   │   ├── network.go       # RouteTable, Route entries
│   │   ├── security.go      # SecurityGroup rules, NetworkACL rules
│   │   └── reachability.go  # VPC Reachability Analyzer
│   ├── config/
│   │   └── config.go        # 기존 (~/.tui-aws/ 경로 변경)
│   ├── store/
│   │   ├── favorites.go     # 기존
│   │   └── history.go       # 기존
│   └── ui/
│       ├── root.go          # RootModel, 탭 전환, StatusBar 탭 표시
│       ├── shared/
│       │   ├── styles.go    # 공용 스타일
│       │   ├── table.go     # 공용 테이블 렌더러
│       │   ├── overlay.go   # 공용 오버레이 렌더러
│       │   └── selector.go  # 공용 선택기
│       ├── tab_ec2/
│       │   ├── model.go     # EC2 서브모델 (기존 로직 추출)
│       │   ├── actions.go   # 액션 메뉴 + Network Path
│       │   └── table.go     # EC2 테이블 컬럼/렌더링
│       ├── tab_vpc/
│       │   ├── model.go     # VPC 목록, 필터, 정렬
│       │   ├── detail.go    # VPC 상세 오버레이 (IGW/NAT/Peering/TGW/Endpoint/EIP)
│       │   └── table.go     # VPC 테이블 컬럼
│       ├── tab_subnet/
│       │   ├── model.go     # Subnet 목록, VPC 필터
│       │   ├── detail.go    # ENI 오버레이
│       │   └── table.go     # Subnet 테이블 컬럼
│       ├── tab_routetable/
│       │   ├── model.go     # Route Table 목록
│       │   ├── detail.go    # Route Entries 오버레이
│       │   └── table.go     # Route Table 테이블 컬럼
│       ├── tab_sg/
│       │   ├── model.go     # SG/NACL 전환, 목록
│       │   ├── detail.go    # 규칙 테이블 오버레이
│       │   └── table.go     # SG/NACL 테이블 컬럼
│       └── tab_troubleshoot/
│           ├── model.go     # 소스/대상 선택 UI
│           ├── checker.go   # 로컬 검증 엔진 (SG+Route+NACL)
│           └── result.go    # 결과 렌더링
```

## 데이터 모델

### VPC 관련 (vpc.go)

```go
type VPC struct {
    VpcID, Name, CIDRBlock, State string
    IsDefault bool
    IGWs      []InternetGateway
    NATs      []NATGateway
    Peerings  []VPCPeering
    TGWs      []TGWAttachment
    Endpoints []VPCEndpoint
    EIPs      []ElasticIP
}

type InternetGateway struct { ID, Name, State string }
type NATGateway struct { ID, Name, SubnetID, PrivateIP, PublicIP, State string }
type VPCPeering struct { ID, Name, RequesterVpcID, AccepterVpcID, State string }
type TGWAttachment struct { ID, TGWID, Name, State string }
type VPCEndpoint struct { ID, Name, ServiceName, Type, State string }
type ElasticIP struct { AllocationID, PublicIP, AssociationID, InstanceID, Name string }
```

### Subnet / ENI (subnet.go)

```go
type Subnet struct {
    SubnetID, Name, VpcID, CIDRBlock, AZ string
    AvailableIPs int
    MapPublicIP  bool
    RouteTableID string
}

type ENI struct {
    ID, Description, SubnetID, PrivateIP, PublicIP, Status string
    AttachedInstanceID string
    SecurityGroups     []string
}
```

### Route Table (network.go)

```go
type RouteTable struct {
    ID, Name, VpcID string
    Subnets         []string
    IsMain          bool
    Routes          []Route
}

type Route struct {
    Destination, Target, State string
}
```

### Security / NACL (security.go)

```go
type SecurityGroup struct {
    ID, Name, Description, VpcID string
    InboundRules, OutboundRules []SGRule
}

type SGRule struct {
    Protocol, PortRange, Source, Description string
}

type NetworkACL struct {
    ID, Name, VpcID string
    IsDefault       bool
    Subnets         []string
    InboundRules, OutboundRules []NACLRule
}

type NACLRule struct {
    RuleNumber int
    Protocol, PortRange, CIDRBlock, Action string
}
```

### Reachability (reachability.go)

```go
type ReachabilityResult struct {
    PathID, AnalysisID, Source, Destination string
    Reachable    bool
    Explanations []string
}
```

## API 호출 전략

| 시점 | 호출 | 이유 |
|------|------|------|
| 탭 전환 시 (lazy) | 해당 탭 목록 API | 불필요한 데이터 미리 안 가져옴 |
| 상세 보기 시 | 연관 리소스 API | IGW/NAT 등은 VPC 상세 진입 시만 |
| Troubleshoot 시 | SG+Route+NACL 일괄 | 분석에 모든 데이터 필요 |
| Reachability | 사용자 명시적 선택 후 | 비용 발생 API — 확인 메시지 표시 |

모든 API는 `ec2.Client` 하나로 호출 (새 서비스 클라이언트 불필요).

## 데이터 캐싱

- 캐시 있음 + 30초 미경과 → 캐시 사용
- 캐시 있음 + 30초 경과 → 백그라운드 리로드 (기존 데이터 표시)
- 캐시 없음 → Loading 표시 + API 호출
- `R` 키: 활성 탭 캐시 삭제 + 리로드
- 프로파일/리전 변경: 전체 캐시 삭제
- SSM 세션 복귀: EC2 탭 캐시만 삭제
- 캐시는 메모리 only (SharedState.Cache)

## 탭별 UI 상세

### 탭 1: EC2 (기존 + 확장)

기존 테이블/기능 그대로. 액션 메뉴에 추가:
- **Network Path**: VPC→Subnet→RouteTable→SG→NACL 요약 오버레이
- **Go to VPC**: VPC 탭으로 전환 + 해당 VPC 자동 선택
- **Go to Subnet**: Subnet 탭으로 전환 + 해당 Subnet 자동 선택

### 탭 2: VPC

테이블: Name, VPC ID, CIDR, State, Default, Subnets 수, IGW 유무

액션 메뉴:
- VPC Details (IGW/NAT/Peering/TGW/Endpoint/EIP 전체 오버레이)
- Subnets in this VPC (Subnet 탭 필터 전환)
- Route Tables (Routes 탭 필터 전환)
- Security Groups (SG 탭 필터 전환)

### 탭 3: Subnets

테이블: Name, Subnet ID, VPC Name, CIDR, AZ, Available IPs, Public 여부

액션 메뉴:
- Subnet Details
- ENIs in this Subnet
- Route Table (연결된 RT 상세)
- Go to VPC

VPC 탭에서 진입 시 VPC 필터 자동 적용.

### 탭 4: Routes

테이블: Name, RT ID, VPC Name, Main 여부, Subnets 수, Routes 수

액션 메뉴:
- Route Entries (Destination → Target → State 테이블)
- Associated Subnets

### 탭 5: SG

기본 모드 — Security Groups:
테이블: Name, SG ID, VPC Name, Inbound 수, Outbound 수, Description

`f` 키로 NACL 모드 전환:
테이블: Name, ACL ID, VPC Name, Default 여부, Subnets 수

SG 액션 메뉴:
- Inbound Rules (Proto/Ports/Source/Description 테이블)
- Outbound Rules
- Referenced by (인스턴스/ENI 목록)

### 탭 6: Troubleshoot

초기 화면: Source/Destination 인스턴스 선택, Protocol/Port 입력

검증 엔진 (로컬):
1. Source SG Outbound 확인
2. Source NACL Outbound 확인
3. Source Route Table 경로 확인
4. Destination NACL Inbound 확인
5. Destination SG Inbound 확인

결과: 각 단계 ✓/✗ 표시, 차단 지점에서 중단, 제안 메시지

Reachability Analyzer (옵션):
- 사용자가 `[R: Reachability Analyzer]` 선택 시
- 비용 경고 확인 후 실행
- `CreateNetworkInsightsPath` → `StartNetworkInsightsAnalysis` → 결과 폴링 → 표시

## 인스턴스→네트워크 드릴다운

EC2 인스턴스에서 양방향 진입 가능:

```
EC2 탭 → Enter → "Network Path"   → VPC/Subnet/Route/SG/NACL 요약 오버레이
EC2 탭 → Enter → "Go to VPC"      → VPC 탭 (해당 VPC 포커스)
EC2 탭 → Enter → "Go to Subnet"   → Subnet 탭 (해당 Subnet 포커스)

VPC 탭 → Enter → "Subnets"        → Subnet 탭 (VPC 필터 적용)
Subnet 탭 → Enter → "Route Table" → Routes 탭 (해당 RT 포커스)
```

## 키 바인딩

### 전역

| 키 | 동작 |
|----|------|
| `1`-`6` | 탭 전환 |
| `Tab`/`Shift+Tab` | 다음/이전 탭 |
| `p` | 프로파일 선택 |
| `r` | 리전 선택 |
| `R` | 활성 탭 새로고침 |
| `q`/`Ctrl+C` | 종료 |

### 탭 내 (공통)

| 키 | 동작 |
|----|------|
| `↑↓`/`jk` | 행 이동 |
| `Enter` | 액션 메뉴 |
| `/` | 검색 |
| `f` | 필터 (SG 탭에서는 SG/NACL 전환) |
| `s`/`S` | 정렬 컬럼/방향 |
| `F` | 즐겨찾기 (EC2 탭만) |
| `Esc` | 오버레이/검색 닫기 |

오버레이 활성 시 숫자키는 포트 입력 등 오버레이 전용. 충돌 없음.

## 터미널 폭 대응

| 폭 | 동작 |
|----|------|
| < 80 | 핵심 컬럼만 (Name, ID, 1-2개), 탭 바 숫자만 `[1][2]...` |
| 80-160 | 기본 컬럼 |
| > 160 | 확장 컬럼 (Description 등) |

## 에러 핸들링

- 탭별 독립: API 실패 시 해당 탭만 에러 표시, 다른 탭 영향 없음
- 권한 부족: `"AccessDenied: ec2:DescribeRouteTables — IAM 정책에 권한을 추가하세요"` 표시
- Reachability Analyzer: 실행 전 비용 경고 확인 메시지

## IAM 권한

기존: `ec2:DescribeInstances`, `ssm:StartSession`, `ssm:DescribeInstanceInformation`, `ec2:DescribeVpcs`, `ec2:DescribeSubnets`

Phase 1 추가: `ec2:DescribeInternetGateways`, `ec2:DescribeNatGateways`, `ec2:DescribeVpcPeeringConnections`, `ec2:DescribeTransitGatewayAttachments`, `ec2:DescribeVpcEndpoints`, `ec2:DescribeAddresses`, `ec2:DescribeNetworkInterfaces`

Phase 2 추가: `ec2:DescribeRouteTables`, `ec2:DescribeSecurityGroups`, `ec2:DescribeNetworkAcls`

Phase 3 추가 (옵션): `ec2:CreateNetworkInsightsPath`, `ec2:StartNetworkInsightsAnalysis`, `ec2:DescribeNetworkInsightsAnalyses`, `ec2:DeleteNetworkInsightsPath`

## 구현 페이즈

### Phase 0: 리팩토링

- go.mod/Makefile/설정 디렉토리 이름 변경
- model.go → root.go + tab_ec2/ 분리
- shared/ 공용 컴포넌트 추출
- TabModel 인터페이스 정의
- 완료 기준: 기존 EC2 기능 동일 동작, 테스트 통과

### Phase 1: VPC + Subnet 탭

- aws/vpc.go, aws/subnet.go 구현
- tab_vpc/, tab_subnet/ 서브모델
- EC2 액션 메뉴에 Go to VPC/Subnet 추가
- 완료 기준: 탭 1-3 동작, 드릴다운 동작

### Phase 2: Routes + SG/NACL 탭

- aws/network.go, aws/security.go 구현
- tab_routetable/, tab_sg/ 서브모델
- EC2 Network Path 오버레이
- 완료 기준: 탭 1-5 동작, Network Path 동작

### Phase 3: Troubleshoot 탭

- tab_troubleshoot/ 서브모델
- 로컬 검증 엔진 (SG+Route+NACL)
- aws/reachability.go (옵션)
- 완료 기준: 탭 6 동작, 로컬 검증 + Reachability Analyzer
