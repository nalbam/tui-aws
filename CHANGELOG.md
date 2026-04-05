# Changelog

[![English](https://img.shields.io/badge/lang-English-blue)](#english)
[![한국어](https://img.shields.io/badge/lang-한국어-red)](#한국어)

---

# English

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-04-05

### Added

- Add 22 integrated tabs for AWS infrastructure management: EC2, ASG, EBS, VPC, Subnet, Routes, SG/NACL, VPCE, TGW, ELB, CloudFront, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda, CloudWatch, IAM, Connectivity Check
- Add SSM Session Manager integration for connecting to EC2 instances without SSH keys
- Add port forwarding via SSM to tunnel local ports to remote instances
- Add network path visualization showing VPC, Subnet, Route Table, Security Group, and NACL in one overlay
- Add local connectivity checker that validates SG + Route + NACL rules between two EC2 instances in 5 steps without calling AWS APIs
- Add AWS Reachability Analyzer integration with cost confirmation prompt
- Add ECS deep dive with hierarchical drill-down: Clusters, Services, Tasks, Containers, CloudWatch Logs, and ECS Exec
- Add EKS K8s integration via direct REST API calls (no kubectl dependency): Pods, Deployments, Services, Nodes, and Pod Logs
- Add instance favorites (persisted to `~/.tui-aws/favorites.json`) and SSM session history tracking
- Add cross-resource navigation between tabs (EC2 to VPC, Subnet, Route Table, Security Group)
- Add AWS profile selector supporting named profiles, SSO, and EC2 instance roles
- Add AWS region selector with all standard regions
- Add interactive ELB target group detail with per-target health status
- Add EBS encryption status column with color-coded indicators
- Add inline search in all tabs filtering by name, ID, and IP
- Add sort and reverse sort across all table views
- Add SG/NACL dual-mode toggle on a single tab
- Add cross-platform setup script for macOS and Linux with automatic prerequisite installation
- Add cross-compilation support for linux/darwin x amd64/arm64

### Fixed

- Fix macOS build compatibility by splitting stdin flush into platform-specific files with Go build tags
- Fix AWS profile credential fallback chain to validate credentials before use and gracefully fall back to instance role
- Fix terminal state corruption after SSM and ECS Exec sessions with `stty sane` and stdin flush
- Fix setup script compatibility on macOS by replacing GNU-specific commands with POSIX alternatives

[Unreleased]: https://github.com/whchoi98/tui-aws/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/whchoi98/tui-aws/releases/tag/v0.1.0

---

# 한국어

이 프로젝트의 모든 주요 변경 사항은 이 파일에 기록됩니다.
이 문서는 [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)를 기반으로 하며,
[Semantic Versioning](https://semver.org/spec/v2.0.0.html)을 따릅니다.

## [Unreleased]

## [0.1.0] - 2026-04-05

### Added

- AWS 인프라 관리를 위한 22개 통합 탭 추가: EC2, ASG, EBS, VPC, Subnet, Routes, SG/NACL, VPCE, TGW, ELB, CloudFront, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda, CloudWatch, IAM, 연결성 검사
- SSH 키 없이 EC2 인스턴스에 접속하는 SSM Session Manager 통합
- 로컬 포트를 원격 인스턴스로 터널링하는 SSM 포트 포워딩 기능 추가
- VPC, Subnet, Route Table, Security Group, NACL을 하나의 오버레이에서 보여주는 네트워크 경로 시각화 추가
- AWS API 호출 없이 두 EC2 인스턴스 간 SG + Route + NACL 규칙을 5단계로 검증하는 로컬 연결성 검사기 추가
- 비용 확인 프롬프트를 포함한 AWS Reachability Analyzer 통합
- 계층적 드릴다운을 지원하는 ECS 딥다이브 추가: Clusters, Services, Tasks, Containers, CloudWatch Logs, ECS Exec
- K8s REST API 직접 호출을 통한 EKS K8s 통합 추가 (kubectl 불필요): Pods, Deployments, Services, Nodes, Pod Logs
- 인스턴스 즐겨찾기 (`~/.tui-aws/favorites.json`에 저장) 및 SSM 세션 이력 추적 추가
- 탭 간 크로스 리소스 탐색 추가 (EC2에서 VPC, Subnet, Route Table, Security Group으로 이동)
- 명명된 프로파일, SSO, EC2 인스턴스 역할을 지원하는 AWS 프로파일 선택기 추가
- 모든 표준 리전을 포함하는 AWS 리전 선택기 추가
- 타겟별 헬스 상태를 보여주는 대화형 ELB 타겟 그룹 상세 추가
- 색상 표시를 포함한 EBS 암호화 상태 컬럼 추가
- 모든 탭에서 이름, ID, IP로 필터링하는 인라인 검색 추가
- 모든 테이블 뷰에서 정렬 및 역순 정렬 추가
- 단일 탭에서 SG/NACL 듀얼 모드 전환 추가
- 사전 요구 사항 자동 설치를 포함하는 macOS 및 Linux용 크로스 플랫폼 설치 스크립트 추가
- linux/darwin x amd64/arm64 크로스 컴파일 지원 추가

### Fixed

- Go 빌드 태그를 사용하여 stdin flush를 플랫폼별 파일로 분리하여 macOS 빌드 호환성 수정
- 자격 증명 사용 전 검증 및 인스턴스 역할로의 우아한 폴백을 위한 AWS 프로파일 자격 증명 폴백 체인 수정
- `stty sane`과 stdin flush를 통한 SSM 및 ECS Exec 세션 후 터미널 상태 손상 수정
- GNU 전용 명령을 POSIX 대체 명령으로 교체하여 macOS에서의 설치 스크립트 호환성 수정

[Unreleased]: https://github.com/whchoi98/tui-aws/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/whchoi98/tui-aws/releases/tag/v0.1.0
