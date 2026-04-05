# Onboarding Guide

<p align="center">
  <kbd><a href="#한국어">한국어</a></kbd> · <kbd><a href="#english">English</a></kbd>
</p>

---

## 한국어

### 사전 요구사항

- **Go 1.25+** — `go version`으로 확인
- **AWS CLI v2** — `aws --version`으로 확인
- **AWS Session Manager Plugin** — SSM 세션 및 ECS Exec에 필요
- **AWS 자격 증명** — `~/.aws/credentials` 또는 `~/.aws/config`에 프로필 설정

### 초기 설정

```bash
# 저장소 클론
git clone <repository-url>
cd tui-aws

# 자동 설치 (Go, AWS CLI, Session Manager Plugin)
bash scripts/setup.sh

# 또는 수동 빌드
make build
```

### 프로젝트 구조 이해

```
main.go              → 진입점, 설정 마이그레이션, TUI 실행
internal/aws/        → AWS SDK 호출 (18개 서비스 클라이언트)
internal/config/     → 사용자 설정 (~/.tui-aws/config.json)
internal/store/      → 즐겨찾기 & 히스토리
internal/ui/root.go  → RootModel: 탭 전환, 글로벌 키
internal/ui/shared/  → TabModel 인터페이스, SharedState, 스타일
internal/ui/tab_*/   → 22개 탭 각각의 구현
```

### 개발 워크플로우

```bash
# 빌드 및 실행
make build && ./tui-aws

# 테스트
make test

# 정적 분석
go vet ./...

# 특정 탭 테스트
go test ./internal/ui/tab_troubleshoot/ -v
```

### 새 탭 추가 방법

1. `internal/ui/shared/tab.go`에 TabID 추가
2. `internal/ui/tab_<name>/` 패키지 생성 (model.go, table.go, detail.go)
3. `internal/ui/root.go`에 탭 등록
4. `CLAUDE.md` 문서 생성

### 핵심 패턴

- **Elm 아키텍처:** Model → Update (메시지 처리) → View (렌더링)
- **지연 로딩:** 탭은 첫 활성화 시 데이터 로드, 30초 캐시 TTL
- **SSM/ECS Exec:** `tea.Exec()`으로 TUI 일시 중단 후 외부 프로세스 실행

---

## English

### Prerequisites

- **Go 1.25+** — verify with `go version`
- **AWS CLI v2** — verify with `aws --version`
- **AWS Session Manager Plugin** — required for SSM sessions and ECS Exec
- **AWS Credentials** — configure profiles in `~/.aws/credentials` or `~/.aws/config`

### Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd tui-aws

# Automated setup (installs Go, AWS CLI, Session Manager Plugin)
bash scripts/setup.sh

# Or build manually
make build
```

### Project Structure

```
main.go              → Entry point, config migration, TUI launch
internal/aws/        → AWS SDK calls (18 service clients)
internal/config/     → User preferences (~/.tui-aws/config.json)
internal/store/      → Favorites & session history
internal/ui/root.go  → RootModel: tab switching, global keys
internal/ui/shared/  → TabModel interface, SharedState, styles
internal/ui/tab_*/   → Implementation of each of the 22 tabs
```

### Development Workflow

```bash
# Build and run
make build && ./tui-aws

# Run tests
make test

# Static analysis
go vet ./...

# Test a specific tab
go test ./internal/ui/tab_troubleshoot/ -v
```

### Adding a New Tab

1. Add TabID to `internal/ui/shared/tab.go`
2. Create `internal/ui/tab_<name>/` package (model.go, table.go, detail.go)
3. Register in `internal/ui/root.go`
4. Create `CLAUDE.md` documentation

### Key Patterns

- **Elm Architecture:** Model -> Update (message handling) -> View (rendering)
- **Lazy Loading:** Tabs fetch data on first activation, 30s cache TTL
- **SSM/ECS Exec:** TUI suspended via `tea.Exec()` to run external processes
