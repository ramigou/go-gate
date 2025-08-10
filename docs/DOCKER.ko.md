# Go-Gate Docker 설정

*다른 언어로 읽기: [English](DOCKER.md), [日本語](DOCKER.ja.md)*

이 가이드는 Docker를 사용하여 Go-Gate 리버스 프록시 서버와 테스트 환경을 실행하는 방법을 설명합니다.

## 빠른 시작

### 사전 요구 사항

- Docker Engine 20.10 이상
- Docker Compose 2.0 이상

### 원클릭 설정

```bash
# 모든 서비스 시작
./docker/start.sh
```

이 명령은:
- Go-Gate 프록시 서버를 빌드
- 3개의 목 백엔드 서버 시작 (API 서버 1, API 서버 2, 웹 서버)
- 8080 포트에서 리버스 프록시 시작
- 헬스 체크 수행
- 서비스 상태 표시

## Docker 구성 요소

### 서비스

| 서비스 | 설명 | 포트 | URL |
|---------|-------------|------|-----|
| `go-gate` | 리버스 프록시 서버 | 8080 | http://localhost:8080 |
| `api-server-1` | 목 API 백엔드 (가중치: 2) | 3001 | http://localhost:3001 |
| `api-server-2` | 목 API 백엔드 (가중치: 1) | 3002 | http://localhost:3002 |
| `web-server` | 목 웹 백엔드 | 4000 | http://localhost:4000 |
| `test-runner` | 자동화된 테스트 컨테이너 | - | - |

### 네트워크

모든 서비스는 격리된 통신을 위해 커스텀 브리지 네트워크 `go-gate-network`에서 실행됩니다.

## 사용법

### 환경 시작

```bash
# 백그라운드에서 모든 서비스 시작
./docker/start.sh

# 또는 docker-compose로 수동 실행
docker-compose up -d
```

### 테스트 실행

```bash
# 자동화된 테스트 실행
./docker/test.sh

# 또는 수동 실행
docker-compose run --rm test-runner
```

### 환경 중지

```bash
# 모든 서비스 중지
./docker/stop.sh

# 또는 수동으로
docker-compose down
```

## 테스트

### 자동화된 테스트

테스트 러너가 포괄적인 테스트를 수행합니다:

```bash
docker-compose run --rm test-runner
```

테스트 내용:
- API 라우트 로드 밸런싱
- 기본 라우트 분산
- 호스트 기반 라우팅
- 서비스 헬스 체크

### 수동 테스트

```bash
# API 라우팅 테스트 (로드 밸런싱)
curl http://localhost:8080/api/users

# 기본 라우팅 테스트
curl http://localhost:8080/default

# 호스트 기반 라우팅 테스트
curl -H "Host: admin.example.com" http://localhost:8080/dashboard
curl -H "Host: www.example.com" http://localhost:8080/home

# POST 요청 테스트
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

## 모니터링

### 로그 보기

```bash
# 모든 서비스
docker-compose logs -f

# 특정 서비스
docker-compose logs -f go-gate
docker-compose logs -f api-server-1

# 프록시 로그만 실시간으로
docker-compose logs -f go-gate | grep -E "(api_server|web_server)"
```

### 서비스 상태

```bash
# 실행 중인 서비스 확인
docker-compose ps

# 서비스 헬스 체크
docker-compose exec go-gate wget --spider -q http://localhost:8080/
```

## 구성

### Docker 구성

Docker 설정은 컨테이너 네트워킹에 최적화된 `configs/docker-config.yaml`을 사용합니다:

- 컨테이너 호스트명 사용 (`api-server-1`, `api-server-2`, `web-server`)
- Docker 네트워크 통신용으로 구성
- 로컬 개발과 동일한 라우팅 규칙

### 환경 변수

환경 변수를 사용하여 구성을 재정의할 수 있습니다:

```bash
# 커스텀 포트
PROXY_PORT=9090 docker-compose up -d

# 커스텀 구성 파일
CONFIG_FILE=configs/production-config.yaml docker-compose up -d
```

## 개발

### 커스텀 이미지 빌드

```bash
# Go-Gate 이미지 빌드
docker-compose build go-gate

# 또는 수동으로
docker build -t go-gate:latest .
```

### 개발 모드

라이브 리로딩을 위한 개발용:

```bash
# 소스 코드 마운트
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## 문제 해결

### 일반적인 문제

**서비스가 시작되지 않음:**
```bash
# Docker 상태 확인
docker info

# 포트 충돌 확인
lsof -i :8080 -i :3001 -i :3002 -i :4000

# 서비스 로그 확인
docker-compose logs
```

**프록시가 502 오류 반환:**
```bash
# 백엔드 서비스 확인
docker-compose ps
docker-compose logs api-server-1 api-server-2 web-server

# 백엔드 서비스 직접 테스트
curl http://localhost:3001/
curl http://localhost:3002/
curl http://localhost:4000/
```

**테스트 실패:**
```bash
# 서비스가 준비될 때까지 대기
sleep 10

# 프록시 헬스 체크
curl -f http://localhost:8080/

# 상세한 출력으로 테스트 실행
docker-compose run --rm test-runner
```

### 정리

```bash
# 모든 컨테이너와 네트워크 제거
docker-compose down

# 컨테이너, 네트워크, 이미지 제거
docker-compose down --rmi local

# 볼륨을 포함한 모든 것 제거
docker-compose down -v --rmi all
```

## 프로덕션 배포

프로덕션 배포를 위해:

1. 프로덕션 구성 사용
2. 적절한 로깅 설정
3. 헬스 모니터링 구성
4. 컨테이너 오케스트레이션 사용 (Kubernetes, Docker Swarm)
5. SSL/TLS 종료 설정
6. 필요시 영구 볼륨 구성

```bash
# 프로덕션 예시
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```