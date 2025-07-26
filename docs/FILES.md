# 📁 파일별 상세 설명서

## 🎯 프로젝트 파일 구조

```text
multiple-example/
├── main.go                    # 🚀 서버 진입점
├── go.mod                     # 📦 Go 모듈 정의
├── README.md                  # 📖 프로젝트 개요
├── docs/                      # 📚 문서 폴더
│   ├── ARCHITECTURE.md       # 🏗️ 아키텍처 설계
│   ├── FILES.md              # 📁 파일별 설명 (현재 문서)
│   └── API.md                # 🔌 API 명세
├── internal/                  # 🔒 내부 패키지
│   ├── models/               # 📊 데이터 모델
│   │   ├── player.go         # 👤 플레이어 구조체
│   │   ├── message.go        # 📨 메시지 타입
│   │   └── game_state.go     # 🎮 게임 상태
│   ├── game/                 # 🎯 게임 로직
│   │   └── game.go           # ⚙️ 게임 엔진
│   └── websocket/            # 🌐 WebSocket 처리
│       └── handler.go        # 🔌 연결 핸들러
└── public/                   # 🌍 정적 파일
    └── index.html            # 🎮 클라이언트 게임
```

---

## 🚀 main.go - 서버 진입점

### 📋 역할

- 애플리케이션의 메인 진입점
- Fiber 웹 서버 초기화 및 설정
- 정적 파일 서빙 설정
- WebSocket 라우팅 설정

### 🔧 주요 구성 요소

```go
func main() {
    // 1. Fiber 앱 초기화
    app := fiber.New()

    // 2. 정적 파일 서빙 설정
    app.Static("/", "./public")

    // 3. 게임 인스턴스 생성
    game := game.NewGame()

    // 4. WebSocket 핸들러 생성
    wsHandler := ws.NewHandler(game)

    // 5. WebSocket 라우팅 설정
    app.Get("/ws", websocket.New(wsHandler.HandleWebSocket))

    // 6. 서버 시작
    app.Listen(":3000")
}
```

### 🎯 핵심 기능

- **정적 파일 서빙**: `./public` 폴더의 파일들을 웹에서 접근 가능
- **WebSocket 엔드포인트**: `/ws` 경로로 WebSocket 연결 수락
- **의존성 주입**: 게임 인스턴스를 WebSocket 핸들러에 주입

---

## 📦 go.mod - Go 모듈 정의

### 📋 역할

- Go 모듈의 메타데이터 정의
- 의존성 관리
- Go 버전 요구사항 명시

### 🔧 주요 내용

```go
module github.com/sangjinsu/websocket-multiplayer

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.52.9
    github.com/gofiber/websocket/v2 v2.2.1
)
```

### 🎯 의존성

- **Fiber v2**: 고성능 웹 프레임워크
- **WebSocket v2**: WebSocket 지원
- **Go 1.21+**: 최신 Go 기능 활용

---

## 📊 internal/models/ - 데이터 모델

### 👤 player.go - 플레이어 구조체

#### 📋 역할

- 플레이어 정보를 담는 데이터 구조 정의
- WebSocket 연결 정보 포함
- JSON 직렬화 지원

#### 🔧 주요 구조체

```go
type Player struct {
    ID          string    `json:"id"`           // 고유 식별자
    PlayerNum   int       `json:"playerNum"`    // 접속 순서
    Name        string    `json:"name"`         // 플레이어 이름
    X           float64   `json:"x"`            // X 좌표
    Y           float64   `json:"y"`            // Y 좌표
    Vx          float64   `json:"vx"`           // X 속도
    Vy          float64   `json:"vy"`           // Y 속도
    Color       string    `json:"color"`        // 플레이어 색상
    JoinedAt    time.Time `json:"joinedAt"`     // 접속 시간
    LastSeen    time.Time `json:"lastSeen"`     // 마지막 활동 시간
    Conn        *websocket.Conn                 // WebSocket 연결
}
```

#### 🎯 핵심 기능

- **위치 정보**: X, Y 좌표로 플레이어 위치 표현
- **속도 정보**: Vx, Vy로 물리 시뮬레이션 지원
- **연결 정보**: WebSocket 연결 객체 포함
- **시간 정보**: 접속 및 활동 시간 추적

### 📨 message.go - 메시지 타입

#### 📋 역할

- WebSocket 통신에 사용되는 메시지 구조 정의
- 다양한 메시지 타입 지원

#### 🔧 주요 구조체

```go
// 기본 메시지 구조
type Message struct {
    Type    MessageType `json:"type"`
    Payload interface{} `json:"payload"`
}

// 플레이어 이동 메시지
type PlayerMove struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}

// 플레이어 입장 메시지
type PlayerJoin struct {
    ID    string `json:"id"`
    Color string `json:"color"`
}

// 플레이어 로그인 메시지
type PlayerLogin struct {
    Name         string  `json:"name"`
    Color        string  `json:"color"`
    LastPosition *struct {
        X float64 `json:"x"`
        Y float64 `json:"y"`
    } `json:"lastPosition"`
}
```

### 🎮 game_state.go - 게임 상태

#### 📋 역할

- 전체 게임 상태를 관리하는 구조체
- 동시성 안전성 보장

#### 🔧 주요 구조체

```go
type GameState struct {
    Mu          sync.RWMutex
    Players     map[string]*Player
    PlayerCount int
}

func NewGameState() *GameState {
    return &GameState{
        Players: make(map[string]*Player),
    }
}
```

#### 🎯 핵심 기능

- **동시성 제어**: `sync.RWMutex`로 스레드 안전성 보장
- **플레이어 관리**: 맵으로 플레이어 저장 및 조회
- **카운터**: 현재 플레이어 수 추적

---

## 🎯 internal/game/ - 게임 로직

### ⚙️ game.go - 게임 엔진

#### 📋 역할

- 게임의 핵심 로직 처리
- 물리 시뮬레이션 엔진
- 플레이어 관리 시스템

#### 🔧 주요 함수들

##### 플레이어 관리

```go
// 플레이어 추가
func (g *Game) AddPlayer(player *models.Player)

// 플레이어 제거
func (g *Game) RemovePlayer(playerID string)

// 플레이어 조회
func (g *Game) GetPlayer(playerID string) *models.Player

// 모든 플레이어 조회
func (g *Game) GetAllPlayers() map[string]*models.Player
```

##### 물리 엔진

```go
// 입력을 속도로 변환
func (g *Game) ApplyInput(playerID, key string) {
    const speed = 2.5
    switch key {
    case "w": p.Vy -= speed
    case "s": p.Vy += speed
    case "a": p.Vx -= speed
    case "d": p.Vx += speed
    }
}

// 물리 시뮬레이션 (60fps)
func (g *Game) Tick() {
    // 1. 속도 적용 및 마찰
    for _, p := range g.State.Players {
        p.X += p.Vx
        p.Y += p.Vy
        p.Vx *= friction
        p.Vy *= friction
    }

    // 2. 경계 처리
    // 3. 플레이어 간 충돌
}
```

#### 🎯 핵심 알고리즘

##### 탄성 충돌 처리

```go
// 플레이어 간 충돌 시 속도 교환
for idA, a := range players {
    for idB, b := range players {
        if idA >= idB { continue }

        dx := b.X - a.X
        dy := b.Y - a.Y
        dist := math.Sqrt(dx*dx + dy*dy)

        if dist < minDistance && dist > 0 {
            // 위치 분리
            overlap := minDistance - dist
            pushX := (dx / dist) * (overlap / 2)
            pushY := (dy / dist) * (overlap / 2)
            a.X -= pushX; a.Y -= pushY
            b.X += pushX; b.Y += pushY

            // 속도 교환
            nx, ny := dx/dist, dy/dist
            va := a.Vx*nx + a.Vy*ny
            vb := b.Vx*nx + b.Vy*ny
            a.Vx += (vb - va) * nx
            a.Vy += (vb - va) * ny
            b.Vx += (va - vb) * nx
            b.Vy += (va - vb) * ny
        }
    }
}
```

---

## 🌐 internal/websocket/ - WebSocket 처리

### 🔌 handler.go - 연결 핸들러

#### 📋 역할

- WebSocket 연결 관리
- 메시지 라우팅 및 처리
- 브로드캐스팅 시스템

#### 🔧 주요 구조체

```go
type Handler struct {
    game        *game.Game
    tickOnce    sync.Once
    lastGameState map[string]*models.Player
}
```

#### 🎯 핵심 함수들

##### 연결 처리

```go
func (h *Handler) HandleWebSocket(c *websocket.Conn) {
    // 1. 고유 플레이어 ID 생성
    playerID := h.game.GenerateID()

    // 2. 임시 플레이어 생성
    player := &models.Player{
        ID:       playerID,
        Conn:     c,
        LastSeen: time.Now(),
    }

    // 3. 물리 tick 루프 시작 (한 번만)
    h.tickOnce.Do(func() {
        go func() {
            ticker := time.NewTicker(16 * time.Millisecond) // 60fps
            defer ticker.Stop()
            for {
                <-ticker.C
                h.game.Tick()
                h.broadcastGameState()
            }
        }()
    })

    // 4. 메시지 처리 루프
    for {
        _, msg, err := c.ReadMessage()
        if err != nil { break }

        var message models.Message
        if err := json.Unmarshal(msg, &message); err != nil {
            continue
        }

        h.handleMessage(player, message)
    }

    // 5. 연결 해제 처리
    h.game.RemovePlayer(playerID)
    h.broadcastPlayerLeave(playerID)
}
```

##### 메시지 처리

```go
func (h *Handler) handleMessage(player *models.Player, message models.Message) {
    switch message.Type {
    case "login":
        // 로그인 처리
        h.handleLogin(player, message)

    case "input":
        // 입력 처리
        if payload, ok := message.Payload.(map[string]interface{}); ok {
            key, _ := payload["key"].(string)
            h.game.ApplyInput(player.ID, key)
        }

    case models.MessageTypeCollision:
        // 충돌 처리
        h.handleCollision(player, message)
    }
}
```

##### 브로드캐스팅

```go
// 효율적인 게임 상태 브로드캐스트
func (h *Handler) broadcastGameState() {
    players := h.game.GetAllPlayers()

    // 변경사항 감지
    currentState, _ := json.Marshal(players)
    if h.lastGameState != nil {
        lastState, _ := json.Marshal(h.lastGameState)
        if string(currentState) == string(lastState) {
            return // 변경사항이 없으면 전송하지 않음
        }
    }

    // 변경사항이 있으면 브로드캐스트
    msg := models.Message{
        Type:    models.MessageTypeGameState,
        Payload: players,
    }

    data, _ := json.Marshal(msg)
    for _, p := range h.game.State.Players {
        if p.Conn != nil {
            _ = p.Conn.WriteMessage(websocket.TextMessage, data)
        }
    }

    h.lastGameState = players
}
```

---

## 🌍 public/index.html - 클라이언트 게임

#### 📋 역할

- 게임 클라이언트 UI
- WebSocket 연결 관리
- Canvas 기반 게임 렌더링
- 사용자 입력 처리

#### 🔧 주요 구성 요소

##### 게임 클래스

```javascript
class MultiplayerGame {
  constructor() {
    this.canvas = document.getElementById("gameCanvas");
    this.ctx = this.canvas.getContext("2d");
    this.socket = null;
    this.players = {};
    this.myId = null;
    this.myColor = null;
    this.playerName = null;
    this.isConnected = false;
    this.isLoggedIn = false;
  }
}
```

##### WebSocket 연결

```javascript
connect() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
        this.isConnected = true;
        this.updateStatus("연결됨! 게임을 시작하세요.");

        // 로그인 메시지 전송
        const loginMessage = {
            type: "login",
            payload: {
                name: this.playerName,
                color: playerData?.color,
                lastPosition: playerData?.lastPosition,
            },
        };
        this.socket.send(JSON.stringify(loginMessage));
    };

    this.socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        this.handleMessage(message);
    };
}
```

##### 입력 처리

```javascript
setupInput() {
    document.addEventListener("keydown", (e) => {
        if (!this.isConnected || !this.isLoggedIn) return;

        const key = e.key.toLowerCase();
        if (["w", "a", "s", "d"].includes(key)) {
            this.socket.send(JSON.stringify({
                type: "input",
                payload: { key, pressed: true },
            }));
        }
    });
}
```

##### 렌더링

```javascript
render() {
    // Canvas 클리어
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    // 3D 배경 그라데이션
    const gradient = this.ctx.createLinearGradient(0, 0, 0, this.canvas.height);
    gradient.addColorStop(0, "rgba(102, 126, 234, 0.1)");
    gradient.addColorStop(0.5, "rgba(118, 75, 162, 0.1)");
    gradient.addColorStop(1, "rgba(0, 0, 0, 0.2)");
    this.ctx.fillStyle = gradient;
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    // 3D 그리드
    this.ctx.strokeStyle = "rgba(255, 255, 255, 0.15)";
    this.ctx.lineWidth = 1;
    // ... 그리드 그리기

    // 플레이어 렌더링
    Object.values(this.players).forEach((player) => {
        this.ctx.save();
        this.ctx.beginPath();
        this.ctx.arc(player.x, player.y, 15, 0, Math.PI * 2);
        this.ctx.fillStyle = player.color || "#fff";
        this.ctx.fill();
        this.ctx.strokeStyle = "#222";
        this.ctx.lineWidth = 2;
        this.ctx.stroke();
        this.ctx.fillStyle = "#fff";
        this.ctx.font = "bold 12px Arial";
        this.ctx.textAlign = "center";
        this.ctx.fillText(player.name || player.id, player.x, player.y + 30);
        this.ctx.restore();
    });
}
```

##### 메시지 처리

```javascript
handleMessage(message) {
    switch (message.type) {
        case "welcome":
            this.myId = message.payload.id;
            this.myColor = message.payload.color;
            this.isLoggedIn = true;
            this.savePlayerData();
            this.updateStatus(`환영합니다! ${message.payload.name} (ID: ${this.myId})`);
            break;

        case "game_state":
            this.players = message.payload;
            this.render();
            this.updatePlayerCount();
            break;

        case "player_join":
            this.players[message.payload.id] = {
                id: message.payload.id,
                playerNum: message.payload.playerNum,
                name: message.payload.name,
                x: message.payload.x,
                y: message.payload.y,
                color: message.payload.color,
            };
            this.updatePlayerCount();
            this.render();
            break;

        case "player_leave":
            delete this.players[message.payload.id];
            this.updatePlayerCount();
            this.render();
            break;
    }
}
```

#### 🎯 핵심 기능

- **실시간 렌더링**: 서버에서 받은 상태로 Canvas 업데이트
- **입력 전송**: WASD 키 입력을 서버로 전송
- **상태 동기화**: 서버의 authoritative 상태에 완전히 의존
- **UI 관리**: 로그인, 플레이어 목록, 연결 상태 표시

---

## 📚 문서 구조

### 📖 README.md

- 프로젝트 개요 및 기능 소개
- 설치 및 실행 방법
- 기술 스택 설명

### 🏗️ docs/ARCHITECTURE.md

- 전체 아키텍처 설계
- 데이터 플로우 설명
- 성능 최적화 전략

### 📁 docs/FILES.md (현재 문서)

- 각 파일별 상세 설명
- 코드 구조 분석
- 핵심 알고리즘 설명

### 🔌 docs/API.md (추후 작성)

- WebSocket API 명세
- 메시지 타입 정의
- 에러 처리 방법

---

이 문서를 통해 프로젝트의 전체 구조와 각 파일의 역할을 이해할 수 있습니다.
