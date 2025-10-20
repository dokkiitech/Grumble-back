grumble-backend/
├─ cmd/
│  ├─ api/
│  │  ├─ main.go                 # HTTPサーバ起動（DI・ルーティング・ミドルウェア）
│  │  └─ wire.go                 # Google Wire等のDI初期化(任意)
│  ├─ batch/
│  │  └─ main.go                 # cronやジョブワーカー起動
│  └─ migrate/
│     └─ main.go                 # DBマイグレーションCLI（任意）
│
├─ api/
│  ├─ openapi/
│  │  ├─ grumble.yml             # OpenAPI定義（単一真実源）
│  │  └─ README.md
│  └─ gen/
│     ├─ http/
│     │  └─ gin/                 # oapi-codegenで生成（Gin版）
│     │     └─ openapi.gen.go
│     └─ client/
│        └─ openapi_client.gen.go# 必要ならクライアント生成
│
├─ internal/
│  ├─ controller/                # 入力境界（アプリ内IF：UseCase呼び出しの窓口）
│  │  ├─ grumble_controller.go
│  │  ├─ vibe_controller.go
│  │  ├─ timeline_controller.go
│  │  ├─ event_controller.go
│  │  └─ auth_controller.go
│  │
│  ├─ handler/                   # トランスポート層（HTTP）← OpenAPIの実体
│  │  └─ http/
│  │     ├─ dto/                 # I/O DTO（OpenAPI typesに薄ラップしてもOK）
│  │     │  ├─ grumble_dto.go
│  │     │  ├─ vibe_dto.go
│  │     │  ├─ timeline_dto.go
│  │     │  ├─ event_dto.go
│  │     │  └─ auth_dto.go
│  │     ├─ middleware/
│  │     │  ├─ request_id.go
│  │     │  ├─ logger.go
│  │     │  ├─ recover.go
│  │     │  ├─ auth_anonymous.go # 匿名認証(デバイス/セッションID)
│  │     │  └─ ratelimit.go
│  │     └─ router/
│  │        ├─ gin/
│  │        │  ├─ router.go      # Gin版
│  │        │  └─ binder.go      # バインド/バリデーション設定
│  │        └─ echo/             # Echo（現在はこちらを採用）
│  │           ├─ router.go
│  │           └─ binder.go
│  │
│  ├─ usecase/                   # アプリケーションユースケース（ビジネスフロー）
│  │  ├─ grumble_post.go         # BE-01 投稿
│  │  ├─ vibe_add.go             # BE-02 共感
│  │  ├─ timeline_get.go         # BE-03 タイムライン
│  │  ├─ purify_check.go         # BE-04 成仏ロジック（しきい値判定）
│  │  ├─ purge_expired.go        # BE-05 24h削除
│  │  ├─ event_manage.go         # BE-06 イベント管理
│  │  └─ auth_anonymous.go       # BE-07 匿名認証
│  │
│  ├─ domain/                    # ドメイン（エンティティ/値オブジェクト/ドメインサービス）
│  │  ├─ entity/
│  │  │  ├─ grumble.go
│  │  │  ├─ vibe.go
│  │  │  ├─ user.go
│  │  │  ├─ poll.go              # 将来機能(現在はファイルのみ作成)
│  │  │  └─ event.go
│  │  ├─ value/
│  │  │  ├─ toxic_level.go
│  │  │  └─ ids.go               # ID生成/検証
│  │  ├─ service/
│  │  │  ├─ purify_service.go    # 成仏の判定/演出トリガー
│  │  │  └─ virtue_service.go    # 徳ポイント加算/ランキング
│  │  └─ error.go                # ドメインエラー型
│  │
│  ├─ repository/                # Port（抽象）: UseCaseが依存するIF
│  │  ├─ grumble_repo.go
│  │  ├─ vibe_repo.go
│  │  ├─ user_repo.go
│  │  ├─ poll_repo.go
│  │  └─ event_repo.go
│  │
│  ├─ infrastructure/            # Adapter（具体）: DB/Cache/Queueなど
│  │  ├─ db/
│  │  │  ├─ postgres/
│  │  │  │  ├─ connection.go
│  │  │  │  ├─ grumble_store.go  # repository.IFの実装
│  │  │  │  ├─ vibe_store.go
│  │  │  │  ├─ user_store.go
│  │  │  │  ├─ poll_store.go
│  │  │  │  └─ event_store.go
│  │  │  └─ sql/
│  │  │     ├─ queries/          # sqlcやプレーンSQLを置く場所（任意）
│  │  │     └─ migrator.go
│  │  ├─ cache/
│  │  │  └─ redis_client.go      # タイムライン/ランキングキャッシュ（任意）
│  │  ├─ queue/
│  │  │  └─ asynq_client.go      # イベント/成仏演出などの非同期処理（任意）
│  │  ├─ id/
│  │  │  └─ uuid_generator.go
│  │  ├─ auth/
│  │  │  └─ anonymous_provider.go# 匿名ユーザーの発行/検証
│  │  └─ clock/
│  │     └─ clock.go             # 時刻抽象（テストしやすさ向上）
│  │
│  ├─ scheduler/                 # バッチ/cron/定期実行
│  │  ├─ cron.go                 # robfig/cron v3など
│  │  ├─ job_purify_check.go     # しきい値超過の成仏フラグ付与
│  │  └─ job_purge_expired.go    # 24h超の削除
│  │
│  ├─ presenter/                 # 出力変換（DTO整形：ハンドラで使う）
│  │  ├─ grumble_presenter.go
│  │  ├─ timeline_presenter.go
│  │  └─ event_presenter.go
│  │
│  └─ shared/
│     ├─ config/
│     │  ├─ config.go            # envやYAMLの読み込み
│     │  └─ default.toml
│     ├─ logger/
│     │  └─ logger.go            # zap/log/slogなど
│     ├─ httpx/
│     │  ├─ error_response.go
│     │  └─ responder.go
│     └─ tracing/
│        └─ otel.go              # OpenTelemetry（任意）
│
├─ migrations/
│  ├─ 0001_init.sql              # AnonymousUser/Grumble/Vibe/Event…
│  └─ 0002_indexes.sql
│
├─ deployments/
│  ├─ docker/
│    ├─ Dockerfile.api
│    ├─ Dockerfile.batch
│    ├─ docker-compose.dev.yml  # Postgres/Redis/Asynqmonなど
│    └─ entrypoint.sh  
│
├─ Makefile
├─ go.mod
├─ go.sum
└─ README.md


