# 本番環境用 Dockerfile
# ビルドステージ
FROM golang:1.24.2-alpine AS builder

# セキュリティアップデートを適用
RUN apk update && apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# Go モジュールファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド（静的バイナリ、最適化）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/main.go

# 実行ステージ
FROM alpine:latest

# セキュリティアップデートとCA証明書をインストール
RUN apk --no-cache add ca-certificates curl tzdata && \
    update-ca-certificates

# 非rootユーザーを作成
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# 作業ディレクトリを設定
WORKDIR /app

# アプリケーションバイナリをコピー
COPY --from=builder /app/main .

# 所有者を変更
RUN chown -R appuser:appgroup /app

# 非rootユーザーに切り替え
USER appuser

# ポートを公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# アプリケーションを実行
CMD ["./main"] 