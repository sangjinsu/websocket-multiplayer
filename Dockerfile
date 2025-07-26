# 🐳 멀티플레이어 게임 Dockerfile

# ========================================
# 1단계: 빌드 스테이지
# ========================================
FROM golang:1.24.5-alpine AS builder

# 메타데이터 설정
LABEL maintainer="sangjinsu"
LABEL description="멀티플레이어 게임 서버"
LABEL version="1.0.0"

# 작업 디렉토리 설정
WORKDIR /app

# 시스템 패키지 업데이트 및 필요한 도구 설치
RUN apk add --no-cache git ca-certificates tzdata

# Go 모듈 파일 복사
COPY go.mod go.sum ./

# 의존성 다운로드 (캐시 최적화)
RUN go mod download

# 소스 코드 복사
COPY . .

# 정적 파일 복사
COPY public/ ./public/

# 빌드 실행 (정적 링킹으로 최적화)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# ========================================
# 2단계: 실행 스테이지
# ========================================
FROM alpine:latest

# 메타데이터 설정
LABEL maintainer="sangjinsu"
LABEL description="멀티플레이어 게임 서버 (최종 이미지)"
LABEL version="1.0.0"

# 보안을 위한 비루트 사용자 생성
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 필요한 패키지 설치
RUN apk --no-cache add ca-certificates tzdata

# 작업 디렉토리 설정
WORKDIR /app

# 빌드 스테이지에서 바이너리 복사
COPY --from=builder /app/main .

# 정적 파일 복사
COPY --from=builder /app/public ./public

# 문서 파일 복사 (선택사항)
COPY --from=builder /app/README.md .
COPY --from=builder /app/docs ./docs

# 소유권 변경
RUN chown -R appuser:appgroup /app

# 비루트 사용자로 전환
USER appuser

# 포트 노출
EXPOSE 3000

# 헬스체크 설정
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# 환경 변수 설정
ENV GIN_MODE=release
ENV PORT=3000

# 컨테이너 시작 명령
CMD ["./main"] 