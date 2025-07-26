[# 🌍 멀티플레이어 게임 (WebSocket + Go)](https://091fade1-0690-450b-b3df-190846ae454a-00-3f0nzfnz89tij.pike.replit.dev/)

실시간 멀티플레이어 게임을 WebSocket과 Go로 구현한 프로젝트입니다. 서버 authoritative physics를 사용하여 안정적이고 동기화된 게임 경험을 제공합니다.

## 🎮 주요 기능

- **실시간 멀티플레이어**: WebSocket을 통한 실시간 통신
- **서버 Authoritative Physics**: 서버에서 모든 물리 연산 처리
- **플레이어 관리**: 로그인, 로그아웃, 재연결 지원
- **물리 시뮬레이션**: 탄성 충돌, 마찰, 경계 처리
- **지속성**: localStorage를 통한 플레이어 정보 저장
- **반응형 UI**: 모던한 디자인과 부드러운 애니메이션

## 🛠 기술 스택

### Backend

- **Go 1.21+**: 서버 로직 및 물리 엔진
- **Fiber v2**: 고성능 웹 프레임워크
- **WebSocket**: 실시간 양방향 통신
- **Goroutines**: 동시성 처리

### Frontend

- **HTML5 Canvas**: 게임 렌더링
- **JavaScript (ES6+)**: 클라이언트 로직
- **CSS3**: 모던 UI/UX 디자인
- **localStorage**: 클라이언트 데이터 저장

### DevOps

- **Docker**: 컨테이너화 및 배포
- **Docker Compose**: 멀티 컨테이너 오케스트레이션
- **Alpine Linux**: 경량 런타임 환경

## 📦 설치 및 실행

### 방법 1: 직접 실행 (개발용)

#### 1. Go 설치

```bash
# Go 1.21 이상 설치 필요
go version
```

#### 2. 프로젝트 클론

```bash
git clone <repository-url>
cd multiple-example
```

#### 3. 의존성 설치

```bash
go mod tidy
```

#### 4. 서버 실행

```bash
go run main.go
```

#### 5. 게임 접속

브라우저에서 `http://localhost:3000` 접속

### 방법 2: Docker 실행 (배포용)

#### 1. Docker 설치

```bash
# Docker 설치 확인
docker --version
docker-compose --version
```

#### 2. 프로젝트 클론

```bash
git clone <repository-url>
cd multiple-example
```

#### 3. Docker 이미지 빌드 및 실행

```bash
# Docker Compose 사용 (권장)
docker-compose up -d

# 또는 직접 Docker 명령어 사용
docker build -t multiplayer-game .
docker run -d -p 3000:3000 --name game-server multiplayer-game
```

#### 4. 게임 접속

브라우저에서 `http://localhost:3000` 접속

#### 5. 컨테이너 관리

```bash
# 로그 확인
docker-compose logs -f game-server

# 서비스 중지
docker-compose down

# 이미지 재빌드
docker-compose up -d --build
```

> 💡 **Docker 사용 시 장점:**
>
> - 환경 의존성 없음
> - 배포 간편함
> - 확장성 우수
> - 보안 강화

## 🎯 게임 조작법

- **WASD 키**: 플레이어 이동
- **로그인**: ID 입력 후 게임 시작
- **멀티플레이**: 여러 브라우저에서 동시 접속 가능

## 🏗 프로젝트 구조

```text
multiple-example/
├── main.go                    # 🚀 서버 진입점
├── go.mod                     # 📦 Go 모듈 정의
├── README.md                  # 📖 프로젝트 개요
├── .gitignore                 # 🚫 Git 무시 파일
├── Dockerfile                 # 🐳 Docker 이미지 정의
├── docker-compose.yml         # 🐳 Docker Compose 설정
├── .dockerignore              # 🐳 Docker 빌드 제외 파일
├── docs/                      # 📚 문서 폴더
│   ├── ARCHITECTURE.md       # 🏗️ 아키텍처 설계
│   ├── FILES.md              # 📁 파일별 설명
│   ├── API.md                # 🔌 API 명세
│   └── DOCKER.md             # 🐳 Docker 배포 가이드
├── public/                    # 🌍 정적 파일
│   └── index.html            # 🎮 클라이언트 게임
└── internal/                  # 🔒 내부 패키지
    ├── models/               # 📊 데이터 모델
    │   ├── player.go         # 👤 플레이어 구조체
    │   ├── message.go        # 📨 메시지 타입
    │   └── game_state.go     # 🎮 게임 상태
    ├── game/                 # 🎯 게임 로직
    │   └── game.go           # ⚙️ 게임 엔진
    └── websocket/            # 🌐 WebSocket 처리
        └── handler.go        # 🔌 연결 핸들러
```

## 🔧 주요 아키텍처

### 서버 Authoritative Physics

- 모든 물리 연산을 서버에서 처리
- 클라이언트는 입력만 전송하고 결과 수신
- 60fps 물리 tick으로 부드러운 시뮬레이션
- 탄성 충돌, 마찰, 경계 처리 구현

### WebSocket 통신

- 실시간 양방향 통신
- JSON 기반 메시지 프로토콜
- 자동 재연결 및 에러 처리
- 효율적인 브로드캐스팅

### 플레이어 관리

- 고유 ID 생성 및 관리
- 로그인/로그아웃 시스템
- 플레이어 정보 지속성
- 순차적 플레이어 번호 할당

## 🚀 성능 최적화

- **메모리 효율성**: 불필요한 로깅 제거
- **네트워크 최적화**: 변경사항이 있을 때만 브로드캐스트
- **렌더링 최적화**: Canvas 기반 효율적 렌더링
- **동시성**: Goroutines를 통한 비동기 처리

## 📝 개발 가이드

### 새로운 기능 추가

1. `internal/models/`에 데이터 구조 정의
2. `internal/game/`에 게임 로직 구현
3. `internal/websocket/`에 메시지 핸들링 추가
4. `public/index.html`에 클라이언트 UI 구현

### 디버깅

- 브라우저 개발자 도구에서 WebSocket 통신 확인
- 서버 로그에서 연결 및 메시지 처리 모니터링

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.

## 🙏 감사의 말

- [Fiber](https://gofiber.io/) - 고성능 Go 웹 프레임워크
- [WebSocket](https://github.com/fasthttp/websocket) - 빠른 WebSocket 구현
- [HTML5 Canvas](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API) - 게임 렌더링

---

**즐거운 게임 되세요! 🎮**
