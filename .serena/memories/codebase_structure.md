# コードベース構造

## ディレクトリツリー

```
grumble-back/
├── cmd/                        # エントリーポイント
│   ├── api/
│   │   └── main.go             # HTTPサーバー起動（DI・ルーティング）
│   ├── batch/
│   │   └── main.go             # cronやジョブワーカー起動
│   └── migrate/
│       └── main.go             # DBマイグレーションCLI（実装中）
│
├── internal/                   # 内部パッケージ（5層構成）
│   ├── controller/             # Controller層
│   │   ├── *_controller.go    # HTTPハンドラー
│   │   ├── *_presenter.go     # レスポンス整形
│   │   ├── dto/                # リクエスト/レスポンスDTO
│   │   ├── middleware/         # 認証、ロギング、リカバリー
│   │   ├── router/             # ルーティング設定（Gin）
│   │   └── httpx/              # HTTPユーティリティ
│   │
│   ├── usecase/                # Usecase層
│   │   ├── grumble_post.go    # BE-01: 投稿
│   │   ├── vibe_add.go        # BE-02: 共感
│   │   ├── timeline_get.go    # BE-03: タイムライン取得
│   │   ├── purify_check.go    # BE-04: 成仏ロジック
│   │   ├── purge_expired.go   # BE-05: 24時間削除
│   │   ├── event_manage.go    # BE-06: イベント管理
│   │   └── auth_anonymous.go  # BE-07: 匿名認証
│   │
│   ├── domain/                 # Domain層
│   │   ├── grumble/            # 愚痴ドメイン
│   │   │   ├── grumble.go     # エンティティ
│   │   │   └── repository.go  # リポジトリIF
│   │   ├── vibe/               # 共感ドメイン
│   │   │   ├── vibe.go
│   │   │   └── repository.go
│   │   ├── user/               # ユーザードメイン
│   │   │   ├── user.go
│   │   │   └── repository.go
│   │   ├── event/              # イベントドメイン
│   │   │   ├── event.go
│   │   │   └── repository.go
│   │   ├── poll/               # 投票ドメイン（将来機能）
│   │   │   ├── poll.go
│   │   │   └── repository.go
│   │   └── shared/             # ドメイン共通
│   │       ├── error.go        # ドメインエラー
│   │       ├── value.go        # 値オブジェクト
│   │       └── service/        # ドメインサービス
│   │           ├── purify_service.go  # 成仏判定
│   │           └── virtue_service.go  # 徳ポイント
│   │
│   ├── infrastructure/         # Infrastructure層（フラット構造）
│   │   ├── grumble_repository.go    # 愚痴リポジトリ実装
│   │   ├── vibe_repository.go       # 共感リポジトリ実装
│   │   ├── user_repository.go       # ユーザーリポジトリ実装
│   │   ├── event_repository.go      # イベントリポジトリ実装
│   │   ├── poll_repository.go       # 投票リポジトリ実装
│   │   ├── cron.go                  # cronジョブ管理
│   │   ├── job_purify_check.go      # 成仏チェックジョブ
│   │   ├── job_purge_expired.go     # 期限切れ削除ジョブ
│   │   └── logger.go                # ロガー実装
│   │
│   ├── api/                    # 生成されたAPIコード
│   │   └── api.gen.go          # oapi-codegenで生成
│   │
│   └── config/
│       └── config.go           # 設定管理（環境変数）
│
├── migrations/                 # DBマイグレーション
│   ├── 0001_init.sql          # 初期スキーマ
│   └── 0002_indexes.sql       # インデックス追加
│
├── oapi-codegen/              # OpenAPIコード生成設定
│   ├── api.yml                # oapi-codegen設定
│   └── README.md              # OpenAPI仕様の説明
│
├── docker/                    # Docker関連
│   ├── Dockerfile             # アプリケーションイメージ
│   └── docker-compose.yml     # PostgreSQL開発環境
│
├── deployments/               # デプロイメント設定
│
├── Makefile                   # ビルド・タスク管理
├── go.mod                     # Go依存関係
├── go.sum                     # 依存関係チェックサム
├── .gitignore                 # Git除外設定
└── README.md                  # プロジェクト説明

```

## 主要ファイルの説明

### エントリーポイント
- `cmd/api/main.go`: HTTPサーバー起動、DI初期化、ルーティング設定
- `cmd/batch/main.go`: バックグラウンドジョブ・cron実行
- `cmd/migrate/main.go`: マイグレーションCLI（実装予定）

### Controller層（HTTPハンドラー）
- リクエストの受け取りとバリデーション
- Usecaseの呼び出し
- レスポンスの整形と返却
- ミドルウェア（認証、ロギング、エラーハンドリング）

### Usecase層（ビジネスフロー）
- ビジネスロジックのオーケストレーション
- 複数のリポジトリやドメインサービスを組み合わせ
- トランザクション管理

### Domain層（コアビジネスロジック）
- エンティティと値オブジェクトの定義
- ドメインロジックの実装
- リポジトリインターフェースの定義
- ドメインサービス（複数のエンティティを横断するロジック）

### Infrastructure層（技術的詳細）
- リポジトリの実装（SQL実行）
- 外部サービスとの連携
- cronジョブの管理
- ロギング実装

### Config層
- 環境変数の読み込み
- アプリケーション設定の管理

## 層間の依存関係

```
Controller → Usecase → Domain ← Infrastructure
                ↓
              Config
```

- **依存の方向**: 外側（Controller）から内側（Domain）へ
- **依存性逆転**: InfrastructureはDomainのインターフェースに依存
- **ドメイン層は独立**: 他のどの層にも依存しない

## OpenAPI仕様の管理

**重要**: このバックエンドは独自のOpenAPI仕様を持たない
- フロントエンドリポジトリ（`github.com/dokkiitech/Grumble`）の `openapi.yaml` を参照
- `make generate` でコード生成時に自動取得
- 生成先: `internal/api/api.gen.go`

## データベース

### テーブル一覧
1. `anonymous_users` - 匿名ユーザー
2. `grumbles` - 愚痴投稿
3. `vibes` - 共感（「わかる…」）
4. `polls` - 投票（将来機能）
5. `poll_votes` - 投票結果（将来機能）
6. `events` - イベント（大怨霊、お焚き上げ）

### マイグレーション
- `migrations/0001_init.sql`: 初期スキーマ作成
- `migrations/0002_indexes.sql`: インデックス追加
- 追加のマイグレーションは連番で作成（`000X_description.sql`）

## 実装状況

### 実装済み
- データベーススキーマ
- ドメインモデルの一部（Grumbleエンティティ）
- oapi-codegen設定
- Dockerによる開発環境

### 実装中・未実装
- ほとんどのController/Usecase/Infrastructure層のコード
- テストコード
- 認証・認可
- ジョブスケジューラー
- ロギング実装

## .gitignoreの対象

- `.claude/` - Claude設定
- `.specify/` - Specify設定
- `.idea/` - JetBrains IDE設定
- `.vscode/` - VS Code設定
- `.DS_Store` - macOS隠しファイル
- `vendor/` - vendorディレクトリ
