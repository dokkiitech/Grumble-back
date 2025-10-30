# コードスタイルと規則

## 言語とコメント
- **コメント**: 日本語で記述
- **変数・関数名**: 英語（CamelCase、PascalCase）

## アーキテクチャ
**シンプルなオニオンアーキテクチャ（5層構成）**

### 1. Controller層 (`internal/controller/`)
- HTTPハンドラーの実装
- リクエスト/レスポンスの変換
- ファイル命名: `*_controller.go`, `*_presenter.go`
- サブディレクトリ: `dto/`, `middleware/`, `router/`, `httpx/`

### 2. Usecase層 (`internal/usecase/`)
- ビジネスロジックのフロー
- 複数のリポジトリやドメインサービスを組み合わせ
- ファイル命名: 機能名.go（例: `grumble_post.go`, `vibe_add.go`）

### 3. Domain層 (`internal/domain/`)
- ドメインモデル（エンティティ、値オブジェクト）
- リポジトリインターフェース
- ドメインサービス
- 構成:
  - `grumble/` - 愚痴ドメイン
  - `vibe/` - 共感ドメイン
  - `user/` - ユーザードメイン
  - `event/` - イベントドメイン
  - `poll/` - 投票ドメイン（将来）
  - `shared/` - ドメイン共通（error.go, value.go, service/）
- 各ドメインディレクトリに `{domain}.go`（エンティティ）と `repository.go`（IF）を配置

### 4. Infrastructure層 (`internal/infrastructure/`)
- リポジトリ実装（フラット構造）
- ジョブ・cron管理
- ロガー実装
- ファイル命名: `*_repository.go`, `job_*.go`, `cron.go`, `logger.go`

### 5. Config層 (`internal/config/`)
- 環境変数管理
- アプリケーション設定

## パッケージ命名
- パッケージ名はディレクトリ名に対応
- 例: `internal/domain/grumble/` → `package model` または `package grumble`

## 構造体とタグ
- 構造体タグを活用: `db:"column_name"`, `json:"field_name"`
- フィールドは公開（PascalCase）
- 例:
```go
type Grumble struct {
    ID         value.GrumbleID  `db:"grumble_id"`
    UserID     value.UserID     `db:"user_id"`
    Content    string           `db:"content"`
}
```

## 値オブジェクト
- カスタム型を `internal/domain/shared/value.go` に定義
- 型安全性を確保（例: `value.GrumbleID`, `value.UserID`, `value.ToxicLevel`）

## エラーハンドリング
- ドメインエラーは `internal/domain/shared/error.go` に定義
- 適切なエラーラッピングを使用

## データベース
- UUIDをPKとして使用（grumbles, anonymous_users）
- BIGSERIALをPKとして使用（vibes, events, polls等）
- タイムスタンプは `TIMESTAMPTZ` 型
- 外部キー制約を明示的に定義
- インデックスを適切に設定

## 命名規則
- ファイル: スネークケース（`grumble_post.go`）
- 型: PascalCase（`GrumbleID`）
- 関数: PascalCase（公開）、camelCase（非公開）
- 定数: PascalCaseまたはALL_CAPS

## インポート
- 標準ライブラリ
- サードパーティライブラリ
- 内部パッケージ
の順にグループ化し、空行で区切る

## リンティング・フォーマット
- 現時点では `.golangci.yml` 等の設定ファイルなし
- `go fmt` での標準フォーマットを使用
- 必要に応じてgolangci-lintを導入予定

## テスト
- 現時点ではテストファイル未実装
- テストファイル命名: `*_test.go`
- テスト関数命名: `Test{Function}` (例: `TestPostGrumble`)
