# 🏗️ 아키텍처 설계 문서

## 📋 개요

이 프로젝트는 **서버 authoritative physics** 패턴을 사용한 실시간 멀티플레이어 게임입니다. 모든 게임 로직과 물리 연산을 서버에서 처리하여 클라이언트 간 동기화 문제를 해결합니다.

## 🎯 설계 원칙

### 1. 서버 Authoritative Physics

- **모든 물리 연산을 서버에서 처리**
- 클라이언트는 입력만 전송하고 결과 수신
- 치팅 방지 및 일관된 게임 상태 보장
- 60fps 물리 tick으로 부드러운 시뮬레이션

### 2. 모듈화 및 관심사 분리

- `models/`: 데이터 구조 정의
- `game/`: 게임 로직 및 물리 엔진
- `websocket/`: 통신 처리
- `public/`: 클라이언트 UI

### 3. 동시성 안전성

- `sync.RWMutex`를 사용한 스레드 안전한 상태 관리
- Goroutines를 통한 비동기 처리
- WebSocket 연결별 독립적인 핸들링

## 🔄 데이터 플로우

```text
클라이언트 입력 → WebSocket → 서버 처리 → 물리 연산 → 브로드캐스트 → 클라이언트 렌더링
```

### 1. 클라이언트 입력

```javascript
// WASD 키 입력
{
  type: "input",
  payload: { key: "w", pressed: true }
}
```

### 2. 서버 처리

```go
// 입력을 속도로 변환
func (g *Game) ApplyInput(playerID, key string) {
  switch key {
  case "w": p.Vy -= speed
  case "s": p.Vy += speed
  // ...
  }
}
```

### 3. 물리 연산 (60fps)

```go
func (g *Game) Tick() {
  // 1. 속도 적용 및 마찰
  // 2. 경계 처리
  // 3. 플레이어 간 충돌
}
```

### 4. 브로드캐스트

```go
// 변경사항이 있을 때만 전송
func (h *Handler) broadcastGameState() {
  if hasChanges {
    sendToAllClients(gameState)
  }
}
```

## 🏛️ 컴포넌트 아키텍처

### Backend Components

#### 1. Main Server (`main.go`)

- **역할**: 애플리케이션 진입점
- **책임**:
  - Fiber 서버 초기화
  - 정적 파일 서빙
  - WebSocket 라우팅 설정

#### 2. Game Engine (`internal/game/game.go`)

- **역할**: 게임 로직 및 물리 엔진
- **책임**:
  - 플레이어 관리 (추가/제거/조회)
  - 물리 연산 (이동/충돌/경계)
  - 게임 상태 관리

#### 3. WebSocket Handler (`internal/websocket/handler.go`)

- **역할**: WebSocket 연결 및 메시지 처리
- **책임**:
  - 연결 수락/해제
  - 메시지 라우팅
  - 브로드캐스팅

#### 4. Data Models (`internal/models/`)

- **역할**: 데이터 구조 정의
- **책임**:
  - 플레이어 정보 구조
  - 메시지 타입 정의
  - 게임 상태 구조

### Frontend Components

#### 1. Game Client (`public/index.html`)

- **역할**: 게임 UI 및 클라이언트 로직
- **책임**:
  - WebSocket 연결 관리
  - 사용자 입력 처리
  - Canvas 렌더링
  - 로컬 상태 관리

## 🔧 핵심 알고리즘

### 1. 물리 시뮬레이션

#### 속도 기반 이동

```go
// 매 틱마다 속도를 위치에 적용
p.X += p.Vx
p.Y += p.Vy
p.Vx *= friction  // 마찰 적용
p.Vy *= friction
```

#### 탄성 충돌

```go
// 플레이어 간 충돌 시 속도 교환
nx, ny := dx/dist, dy/dist  // 법선 벡터
va := a.Vx*nx + a.Vy*ny     // 속도 성분
vb := b.Vx*nx + b.Vy*ny
a.Vx += (vb - va) * nx      // 속도 교환
a.Vy += (vb - va) * ny
```

### 2. 동시성 제어

#### 읽기/쓰기 뮤텍스

```go
type GameState struct {
    Mu      sync.RWMutex
    Players map[string]*Player
}

// 읽기 작업
func (g *Game) GetPlayer(id string) *Player {
    g.State.Mu.RLock()
    defer g.State.Mu.RUnlock()
    return g.State.Players[id]
}

// 쓰기 작업
func (g *Game) AddPlayer(player *Player) {
    g.State.Mu.Lock()
    defer g.State.Mu.Unlock()
    g.State.Players[player.ID] = player
}
```

### 3. 효율적인 브로드캐스팅

#### 변경사항 감지

```go
// 이전 상태와 비교하여 변경사항이 있을 때만 전송
currentState, _ := json.Marshal(players)
if string(currentState) != string(lastState) {
    broadcastToAllClients(currentState)
}
```

## 🚀 성능 최적화

### 1. 메모리 최적화

- 불필요한 로깅 제거
- 객체 재사용
- 효율적인 데이터 구조

### 2. 네트워크 최적화

- 변경사항이 있을 때만 브로드캐스트
- JSON 직렬화 최적화
- 압축 전송 (필요시)

### 3. 렌더링 최적화

- Canvas 기반 효율적 렌더링
- 불필요한 DOM 조작 최소화
- 애니메이션 최적화

## 🔒 보안 고려사항

### 1. 입력 검증

- 클라이언트 입력의 유효성 검사
- 비정상적인 값 필터링

### 2. 서버 Authoritative

- 모든 게임 로직을 서버에서 처리
- 클라이언트 조작 방지

### 3. 연결 관리

- 연결 상태 모니터링
- 비정상 연결 자동 해제

## 📈 확장성

### 1. 수평 확장

- Redis를 통한 상태 공유
- 로드 밸런서를 통한 분산 처리

### 2. 기능 확장

- 채팅 시스템
- 게임 룸 시스템
- 점수 시스템
- 아이템 시스템

### 3. 성능 확장

- 데이터베이스 연동
- 캐싱 시스템
- CDN 활용

## 🧪 테스트 전략

### 1. 단위 테스트

- 각 컴포넌트별 독립적 테스트
- 물리 엔진 로직 검증

### 2. 통합 테스트

- WebSocket 통신 테스트
- 멀티플레이어 시나리오 테스트

### 3. 성능 테스트

- 동시 접속자 수 테스트
- 네트워크 지연 시뮬레이션

---

이 아키텍처는 확장 가능하고 유지보수가 용이한 멀티플레이어 게임 시스템을 제공합니다.
