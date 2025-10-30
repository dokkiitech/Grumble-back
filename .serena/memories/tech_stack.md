# テックスタック

## 言語・ランタイム
- **Go**: 1.23.0+ (toolchain: 1.24.5)
- **プラットフォーム**: Darwin (macOS)

## Webフレームワーク
- **Gin**: v1.11.0 - HTTPウェブフレームワーク

## データベース
- **PostgreSQL**: 16-alpine
  - ホスト: localhost (docker-compose使用時)
  - ポート: 5432
  - ユーザー: grumble
  - パスワード: grumble
  - データベース名: grumble
  - 接続文字列: `postgres://grumble:grumble@localhost:5432/grumble?sslmode=disable`

## コード生成
- **oapi-codegen**: OpenAPI仕様からGoコード生成
  - 設定ファイル: `oapi-codegen/api.yml`
  - 出力先: `internal/api/api.gen.go`
  - Ginサーバー、モデル、strict-serverを生成

## 主要依存パッケージ
- `github.com/gin-gonic/gin` - HTTPフレームワーク
- `github.com/oapi-codegen/nullable` - Nullableサポート
- `github.com/oapi-codegen/runtime` - oapi-codegenランタイム
- `github.com/google/uuid` - UUID生成

## インフラ・ツール
- **Docker Compose**: ローカル開発環境
- **Make**: ビルド・タスク管理

## OpenAPI仕様の管理
**重要**: バックエンドは独自のOpenAPI仕様を持たない
- フロントエンドリポジトリ（`git@github.com:dokkiitech/Grumble.git`）の `openapi.yaml` を参照
- `make generate` で自動取得してコード生成
- ローカルに `../Grumble/openapi.yaml` がある場合はそちらを優先
- なければGitHubから自動ダウンロード
