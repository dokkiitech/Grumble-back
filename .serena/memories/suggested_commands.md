# 推奨コマンド

## 初期セットアップ

### 依存関係のインストール
```bash
make init
```
- Go依存関係のダウンロードと整理
- oapi-codegenのインストール（未インストールの場合）

### OpenAPIコード生成
```bash
make generate
```
- フロントエンドリポジトリから `openapi.yaml` を取得
- `internal/api/api.gen.go` を生成
- ローカルに `../Grumble/openapi.yaml` があればそれを使用
- なければGitHubから自動ダウンロード

### データベース起動
```bash
docker compose -f docker/docker-compose.yml up -d
```
- PostgreSQL 16-alpineコンテナを起動
- ポート5432で公開

### 環境変数設定
```bash
export DATABASE_URL="postgres://grumble:grumble@localhost:5432/grumble?sslmode=disable"
export GRUMBLE_HTTP_ADDR=":8080"
```

## 開発コマンド

### APIサーバー起動
```bash
go run ./cmd/api
```
- HTTPサーバーを起動（デフォルト: :8080）
- マイグレーション自動実行

### バッチ/cronワーカー起動
```bash
go run ./cmd/batch
```
- 成仏チェック、24時間削除などのジョブを実行

### マイグレーション実行（個別）
```bash
go run ./cmd/migrate
```
※ 現在実装中

## テスト

### 全テスト実行
```bash
make test
# または
go test -v ./...
```

### 特定パッケージのテスト
```bash
go test -v ./internal/usecase/...
```

## クリーンアップ

### 生成ファイルの削除
```bash
make clean
```
- `internal/api/api.gen.go` を削除

## Vendoring

### vendorディレクトリ作成
```bash
make vendor
# または
go mod vendor
```

## データベース操作

### データベース停止
```bash
docker compose -f docker/docker-compose.yml down
```

### データベース削除（ボリューム含む）
```bash
docker compose -f docker/docker-compose.yml down -v
```

### データベースログ確認
```bash
docker compose -f docker/docker-compose.yml logs -f db
```

## Go標準コマンド

### 依存関係の更新
```bash
go mod tidy
```

### フォーマット
```bash
go fmt ./...
```

### ビルド
```bash
go build -o bin/api ./cmd/api
go build -o bin/batch ./cmd/batch
```

## Darwin (macOS) 固有のユーティリティ

### ファイル検索
```bash
find . -name "*.go" -type f
```

### パターン検索
```bash
grep -r "pattern" ./internal
```

### ディレクトリ一覧
```bash
ls -la
```

### ポート使用確認
```bash
lsof -i :8080
```

## Git操作

### ブランチ確認
```bash
git branch
```

### 現在のブランチ
- fix/Change-installation-location

### メインブランチ
- main（PR作成時はこちらへ）

### ステータス確認
```bash
git status
```

### コミット
```bash
git add .
git commit -m "message"
```

### プッシュ
```bash
git push origin [branch-name]
```
