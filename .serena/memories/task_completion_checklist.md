# タスク完了時のチェックリスト

## コード変更後の必須手順

### 1. OpenAPI仕様が変更された場合
```bash
make generate
```
- フロントエンドの `openapi.yaml` が更新された場合に実行
- `internal/api/api.gen.go` が再生成される

### 2. フォーマット
```bash
go fmt ./...
```
- 全Goファイルを標準フォーマットに整形
- コミット前に必ず実行

### 3. 依存関係の整理
```bash
go mod tidy
```
- 新しいパッケージを追加/削除した場合に実行
- `go.mod` と `go.sum` を更新

### 4. テスト実行
```bash
make test
# または
go test -v ./...
```
- 全てのテストがパスすることを確認
- ※ 現在はテストファイルが未実装のため、将来的に実施

### 5. ビルド確認
```bash
go build -o bin/api ./cmd/api
```
- ビルドエラーがないことを確認
- または `go run ./cmd/api` で起動確認

## データベース関連の変更時

### マイグレーションファイル追加時
1. `migrations/` ディレクトリに新しいSQLファイルを追加
2. ファイル命名: `000X_description.sql`（連番を振る）
3. APIサーバーを再起動してマイグレーションを実行

### スキーマ変更後の確認
```bash
# データベースに接続して確認
docker compose -f docker/docker-compose.yml exec db psql -U grumble -d grumble
# または
psql postgres://grumble:grumble@localhost:5432/grumble
```

## コミット前のチェック

### 必須確認項目
- [ ] `go fmt ./...` を実行した
- [ ] `go mod tidy` を実行した（依存関係変更時）
- [ ] ビルドエラーがない
- [ ] テストがパスする（実装後）
- [ ] 生成ファイル（`internal/api/api.gen.go`）がコミット対象に含まれている（OpenAPI変更時）
- [ ] `.gitignore` で除外すべきファイルを追加していない
- [ ] 日本語コメントが適切に記述されている
- [ ] エラーハンドリングが適切に実装されている

## レビュー依頼前のチェック

### コード品質
- [ ] オニオンアーキテクチャに従っている
- [ ] 適切な層に責務が配置されている
- [ ] ドメインモデルが適切に設計されている
- [ ] リポジトリインターフェースと実装が分離されている
- [ ] 値オブジェクトを適切に使用している
- [ ] エラーハンドリングが適切

### ドキュメント
- [ ] 複雑なロジックにコメントを追加
- [ ] README.mdの更新（必要な場合）
- [ ] 新機能の説明を追加（必要な場合）

## 将来的に追加予定

### リンティング（golangci-lintの導入後）
```bash
golangci-lint run
```

### 静的解析
```bash
go vet ./...
```

### カバレッジ計測
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### ベンチマーク
```bash
go test -bench=. ./...
```
