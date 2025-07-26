# 🐳 Docker 배포 가이드

## 📋 개요

이 문서는 멀티플레이어 게임을 Docker를 사용하여 배포하는 방법을 설명합니다.

## 🏗️ Docker 아키텍처

### 멀티스테이지 빌드

- **1단계 (Builder)**: Go 소스 코드 컴파일
- **2단계 (Runtime)**: 최소한의 런타임 환경

### 보안 고려사항

- 비루트 사용자 실행
- 최소 권한 원칙
- 정적 링킹으로 의존성 최소화

## 🚀 빠른 시작

### 1. Docker 이미지 빌드

```bash
# 기본 빌드
docker build -t multiplayer-game .

# 태그 지정
docker build -t multiplayer-game:v1.0.0 .

# 플랫폼 지정 (ARM64 등)
docker build --platform linux/amd64 -t multiplayer-game .
```

### 2. 컨테이너 실행

```bash
# 기본 실행
docker run -d -p 3000:3000 --name game-server multiplayer-game

# 백그라운드 실행 + 자동 재시작
docker run -d -p 3000:3000 --restart unless-stopped --name game-server multiplayer-game

# 환경 변수 설정
docker run -d -p 3000:3000 \
  -e PORT=3000 \
  -e GIN_MODE=release \
  --name game-server multiplayer-game
```

### 3. Docker Compose 사용

```bash
# 서비스 시작
docker-compose up -d

# 로그 확인
docker-compose logs -f game-server

# 서비스 중지
docker-compose down

# 이미지 재빌드
docker-compose up -d --build
```

## 🔧 고급 설정

### 1. 환경 변수

```bash
# 환경 변수 파일 생성
cat > .env << EOF
PORT=3000
GIN_MODE=release
MAX_PLAYERS=100
TICK_RATE=60
EOF

# Docker Compose에서 사용
docker-compose --env-file .env up -d
```

### 2. 볼륨 마운트

```bash
# 로그 디렉토리 마운트
docker run -d -p 3000:3000 \
  -v $(pwd)/logs:/app/logs \
  --name game-server multiplayer-game

# 설정 파일 마운트
docker run -d -p 3000:3000 \
  -v $(pwd)/config:/app/config \
  --name game-server multiplayer-game
```

### 3. 네트워크 설정

```bash
# 커스텀 네트워크 생성
docker network create game-network

# 네트워크 지정하여 실행
docker run -d -p 3000:3000 \
  --network game-network \
  --name game-server multiplayer-game
```

## 📊 모니터링

### 1. 컨테이너 상태 확인

```bash
# 실행 중인 컨테이너 확인
docker ps

# 컨테이너 상세 정보
docker inspect game-server

# 리소스 사용량 확인
docker stats game-server
```

### 2. 로그 확인

```bash
# 실시간 로그
docker logs -f game-server

# 최근 로그
docker logs --tail 100 game-server

# 특정 시간 이후 로그
docker logs --since "2025-07-26T23:00:00" game-server
```

### 3. 헬스체크

```bash
# 헬스체크 상태 확인
docker inspect --format='{{.State.Health.Status}}' game-server

# 헬스체크 로그
docker inspect --format='{{range .State.Health.Log}}{{.Output}}{{end}}' game-server
```

## 🔄 배포 전략

### 1. 개발 환경

```yaml
# docker-compose.dev.yml
version: "3.8"
services:
  game-server:
    build: .
    ports:
      - "3000:3000"
    volumes:
      - .:/app # 소스 코드 마운트
    environment:
      - GIN_MODE=debug
    command: ["go", "run", "main.go"]
```

### 2. 프로덕션 환경

```yaml
# docker-compose.prod.yml
version: "3.8"
services:
  game-server:
    image: multiplayer-game:latest
    ports:
      - "3000:3000"
    restart: unless-stopped
    environment:
      - GIN_MODE=release
    volumes:
      - ./logs:/app/logs
```

### 3. 스케일링

```bash
# 여러 인스턴스 실행
docker-compose up -d --scale game-server=3

# 로드 밸런서 설정 (예: nginx)
docker run -d -p 80:80 \
  -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf \
  nginx:alpine
```

## 🛠️ 문제 해결

### 1. 빌드 문제

```bash
# 캐시 없이 빌드
docker build --no-cache -t multiplayer-game .

# 빌드 과정 상세 출력
docker build --progress=plain -t multiplayer-game .
```

### 2. 실행 문제

```bash
# 컨테이너 내부 접속
docker exec -it game-server sh

# 포트 확인
docker port game-server

# 네트워크 연결 확인
docker network inspect game-network
```

### 3. 성능 문제

```bash
# 리소스 제한 설정
docker run -d -p 3000:3000 \
  --memory=512m \
  --cpus=1.0 \
  --name game-server multiplayer-game

# 성능 모니터링
docker stats --no-stream game-server
```

## 🔒 보안 설정

### 1. 사용자 권한

```dockerfile
# Dockerfile에서 이미 설정됨
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser
```

### 2. 네트워크 보안

```bash
# 내부 네트워크만 사용
docker run -d \
  --network game-network \
  --expose 3000 \
  --name game-server multiplayer-game
```

### 3. 리소스 제한

```bash
# 메모리 및 CPU 제한
docker run -d -p 3000:3000 \
  --memory=1g \
  --cpus=2.0 \
  --name game-server multiplayer-game
```

## 📈 성능 최적화

### 1. 이미지 크기 최적화

```dockerfile
# 멀티스테이지 빌드 사용
# 불필요한 파일 제거
# 정적 링킹 사용
```

### 2. 빌드 캐시 최적화

```dockerfile
# 의존성 다운로드를 먼저 수행
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드는 나중에 복사
COPY . .
```

### 3. 런타임 최적화

```bash
# 컨테이너 리소스 제한
# 적절한 헬스체크 설정
# 로그 로테이션 설정
```

## 🚀 배포 자동화

### 1. GitHub Actions

```yaml
# .github/workflows/docker.yml
name: Docker Build and Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build Docker image
        run: docker build -t multiplayer-game .
      - name: Push to registry
        run: |
          docker tag multiplayer-game ${{ secrets.REGISTRY }}/multiplayer-game
          docker push ${{ secrets.REGISTRY }}/multiplayer-game
```

### 2. CI/CD 파이프라인

```bash
# 자동 배포 스크립트
#!/bin/bash
set -e

# 이미지 빌드
docker build -t multiplayer-game .

# 기존 컨테이너 중지
docker stop game-server || true
docker rm game-server || true

# 새 컨테이너 시작
docker run -d -p 3000:3000 \
  --restart unless-stopped \
  --name game-server multiplayer-game

echo "배포 완료!"
```

## 📝 체크리스트

### 배포 전 확인사항

- [ ] Docker 이미지 빌드 성공
- [ ] 컨테이너 정상 실행
- [ ] 포트 접근 가능
- [ ] 헬스체크 통과
- [ ] 로그 정상 출력
- [ ] 환경 변수 설정 완료
- [ ] 볼륨 마운트 확인
- [ ] 네트워크 연결 확인

### 운영 중 모니터링

- [ ] 컨테이너 상태 확인
- [ ] 리소스 사용량 모니터링
- [ ] 로그 분석
- [ ] 성능 지표 확인
- [ ] 보안 업데이트 적용

---

이 가이드를 통해 Docker를 사용한 안전하고 효율적인 배포가 가능합니다.
