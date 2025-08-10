# Go-Gate 리버스 프록시 테스트

*다른 언어로 읽기: [English](README.md), [日本語](README.ja.md)*

이 디렉토리는 Go-Gate L7 리버스 프록시 서버를 위한 테스트 도구와 스크립트를 포함합니다.

## 빠른 시작

### 1. 목 백엔드 서버 시작

```bash
# 터미널 1: 목 백엔드 서버 시작
cd test/
python3 mock-servers.py
```

이 명령은 세 개의 목 HTTP 서버를 시작합니다:
- `localhost:3001` - API 서버 1 (가중치: 2)
- `localhost:3002` - API 서버 2 (가중치: 1)  
- `localhost:4000` - 웹 서버 (가중치: 1)

### 2. 프록시 서버 시작

```bash
# 터미널 2: 리버스 프록시 시작
./go-gate -config configs/test-config.yaml
```

프록시는 `localhost:8080`에서 시작되고 요청을 목 서버로 라우팅합니다.

### 3. 테스트 실행

```bash
# 터미널 3: 자동화된 테스트 실행
cd test/
./test-requests.sh
```

## 테스트 파일

### `mock-servers.py`
테스트용 세 개의 HTTP 서버를 생성하는 Python 스크립트:
- 서버 정보, 요청 세부사항, 타임스탬프가 포함된 JSON 응답 반환
- GET 및 POST 요청 모두 처리
- 헬스 체크용 `/health` 엔드포인트 포함
- 모든 수신 요청 로깅

### `test-requests.sh`
다양한 프록시 시나리오를 테스트하는 셸 스크립트:
- 경로 기반 라우팅 (`/api/*`)
- 호스트 기반 라우팅 (`admin.example.com`, `*.example.com`)
- 로드 밸런싱 검증
- 기본 라우팅 폴백
- POST 요청 처리
- 헬스 체크 엔드포인트

### `test-config.yaml`
테스트용 간소화된 구성 (헬스 체크 설정 제외).

## 수동 테스트

### 기본 요청

```bash
# 기본 라우팅 테스트 (모든 서버에 로드 밸런싱)
curl http://localhost:8080/

# API 라우팅 테스트 (API 서버들 간 로드 밸런싱)
curl http://localhost:8080/api/users

# POST 요청 테스트
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

### 호스트 기반 라우팅

먼저, `/etc/hosts` 파일에 테스트 도메인을 추가합니다:
```bash
# /etc/hosts에 다음 줄들을 추가
127.0.0.1 admin.example.com
127.0.0.1 www.example.com
```

그런 다음 테스트:
```bash
# API 서버 1로만 라우팅되어야 함
curl -H "Host: admin.example.com" http://localhost:8080/dashboard

# 웹 서버로만 라우팅되어야 함
curl -H "Host: www.example.com" http://localhost:8080/home
```

## 예상 동작

### 경로 기반 라우팅
- `/api/*`로의 요청은 API 서버 1 (가중치 2)과 API 서버 2 (가중치 1) 간에 로드 밸런싱됨
- API 서버 1이 약 67%의 요청을, API 서버 2가 약 33%의 요청을 받아야 함

### 호스트 기반 라우팅
- `admin.example.com` 요청은 API 서버 1로만 이동
- `*.example.com` 요청 (`www.example.com`과 같은)은 웹 서버로만 이동

### 기본 라우팅
- 다른 모든 요청은 세 서버 모두에 로드 밸런싱됨
- 분산: API 서버 1 (50%), API 서버 2 (25%), 웹 서버 (25%)

## 문제 해결

### 프록시가 시작되지 않음
- 포트 8080이 사용 가능한지 확인: `lsof -i :8080`
- 구성 문법 검증: `./go-gate -config configs/test-config.yaml`

### 목 서버가 시작되지 않음
- 포트 3001, 3002, 4000이 사용 가능한지 확인
- Python 3이 설치되어 있는지 확인: `python3 --version`

### 요청 실패
- 모든 서버가 실행 중인지 확인: `ps aux | grep -E "(go-gate|python)"`
- 오류 메시지에 대한 프록시 로그 확인
- 목 서버를 직접 테스트: `curl http://localhost:3001/`

### 호스트 기반 라우팅이 작동하지 않음
- `/etc/hosts` 항목이 올바르게 추가되었는지 확인
- 테스트에 `curl -H "Host: domain.com"` 사용
- DNS 해석 확인: `nslookup admin.example.com`

## 고급 테스트

### 로드 테스트
```bash
# 로드 분산 테스트 ('ab' - Apache Bench 필요)
ab -n 100 -c 10 http://localhost:8080/api/test

# 또는 반복문에서 curl 사용
for i in {1..20}; do
  curl -s http://localhost:8080/api/test | grep '"server"'
done | sort | uniq -c
```

### 오류 시나리오
```bash
# 하나의 백엔드 서버를 중지하고 페일오버 테스트
# API 서버 2 중지 (mock-servers.py 터미널에서 Ctrl+C, 그다음 포트 3002 없이 재시작)
curl http://localhost:8080/api/test  # API 서버 1로 계속 작동해야 함

# 모든 백엔드가 다운된 상태로 테스트
# mock-servers.py를 완전히 중지
curl http://localhost:8080/  # 502 Bad Gateway를 반환해야 함
```

## 모니터링

### 실시간 로그
요청 라우팅을 보기 위해 프록시 로그 확인:
```bash
./go-gate -config configs/test-config.yaml | grep -E "(api_server|web_server)"
```

### 요청 분산 분석
```bash
# 어떤 서버가 요청을 처리하고 있는지 분석
./test-requests.sh 2>&1 | grep "Server=" | sort | uniq -c
```