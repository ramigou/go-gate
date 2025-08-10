# Go-Gate L7 리버스 프록시

*다른 언어로 읽기: [English](README.md), [日本語](README.ja.md)*

Go로 구축된 경량 Layer 7 리버스 프록시 서버로, 호스트 기반 및 경로 기반 라우팅과 로드 밸런싱을 지원합니다.

## 기능

- **L7 라우팅**: 호스트, 경로 또는 둘 다를 기준으로 요청 라우팅
- **로드 밸런싱**: 여러 업스트림에 걸친 가중 라운드 로빈 로드 밸런싱
- **헬스 체크**: 기본 헬스 체크 구성 (구현을 위한 구조 준비)
- **요청 로깅**: 응답 시간과 상태를 포함한 HTTP 접근 로그
- **유연한 구성**: 검증 기능이 포함된 YAML 기반 구성

## 프로젝트 구조

```
go-gate/
├── cmd/server/main.go          # 애플리케이션 진입점
├── internal/
│   ├── config/config.go        # 구성 관리
│   ├── proxy/server.go         # 메인 프록시 서버 로직
│   ├── router/router.go        # 요청 라우팅 및 업스트림 선택
│   └── middleware/logging.go   # HTTP 미들웨어 (로깅)
├── configs/config.yaml         # 예제 구성
└── docs/README.md             # 이 문서
```

## 시작하기

### 사전 요구 사항

- Go 1.24.6 이상
- 리버스 프록시와 HTTP에 대한 기본 이해

### 서버 실행

1. **애플리케이션 빌드:**
   ```bash
   go build -o go-gate cmd/server/main.go
   ```

2. **기본 구성으로 실행:**
   ```bash
   ./go-gate
   ```

3. **사용자 정의 구성으로 실행:**
   ```bash
   ./go-gate -config /path/to/your/config.yaml
   ```

### 구성

서버는 YAML 구성 파일을 사용합니다. 예제는 `configs/config.yaml`을 참조하세요.

#### 구성 구조:

- **server**: 서버 설정 (포트)
- **upstreams**: 가중치 및 헬스 체크 설정을 포함한 백엔드 서버 정의
- **routes**: 호스트/경로 매칭을 포함한 라우팅 규칙

#### 라우트 매칭:

- **호스트 매칭**: 정확한 매치와 와일드카드 (`*.example.com`) 지원
- **경로 매칭**: 정확한 매치와 접두사 매칭 (`/api/*`) 지원
- **로드 밸런싱**: 여러 업스트림 간 가중 분산

### 사용 예제

제공된 구성으로 프록시는 다음과 같이 동작합니다:

1. `/api/*` 요청을 API 서버(api_server_1, api_server_2)로 라우팅
2. `admin.example.com` 요청을 api_server_1로 라우팅
3. `*.example.com` 요청을 web_server로 라우팅
4. 기타 요청은 모든 서버로 폴백

### 개발 명령어

- **빌드**: `go build cmd/server/main.go`
- **실행**: `go run cmd/server/main.go`
- **테스트**: `go test ./...`
- **포맷**: `go fmt ./...`
- **검사**: `go vet ./...`
- **모듈 정리**: `go mod tidy`

### 향후 개선 사항

1. **헬스 체킹**: 실제 헬스 체크 모니터링 구현
2. **메트릭**: 모니터링을 위한 Prometheus 메트릭 추가
3. **TLS 지원**: HTTPS/TLS 종료 추가
4. **요청 제한**: 요청 속도 제한 구현
5. **서킷 브레이커**: 업스트림 장애를 위한 서킷 브레이커 패턴 추가
6. **WebSocket 지원**: WebSocket 프록시 기능 추가