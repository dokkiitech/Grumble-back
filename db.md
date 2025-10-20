# 技育ハッカソンDB設計 - Grumble (PostgreSQL)

## 概要
愚痴投稿アプリ「Grumble」のデータベース設計。24時間で消える投稿機能、共感システム、匿名ユーザー管理、イベント機能を含む。

## テーブル設計

### 1. 投稿テーブル (Grumble)
- grumble_id: UUID/SERIAL (PK) - 投稿の一意識別子
- user_id: UUID (FK→AnonymousUser) - 投稿者の匿名ID
- content: TEXT (NOT NULL, 280文字制限) - 愚痴の本文
- toxic_level: INT (NOT NULL, 1-5) - 投稿者の自己申告毒レベル
- vibe_count: INT (NOT NULL, DEFAULT 0) - 「わかる…」の総数（キャッシュ）
- is_purified: BOOLEAN (NOT NULL, DEFAULT FALSE) - 成仏フラグ（Trueでタイムラインから消える）
- posted_at: TIMESTAMP WITH TIME ZONE (NOT NULL) - 投稿時刻（24時間削除の基準）
- expires_at: TIMESTAMP WITH TIME ZONE (NOT NULL) - 投稿後24時間後の時刻
- is_event_grumble: BOOLEAN (NOT NULL, DEFAULT FALSE) - イベント投稿フラグ

**インデックス**: posted_at, expires_at, is_purified

### 2. 共感テーブル (Vibe)
- vibe_id: SERIAL (PK) - 共感履歴の一意識別子
- grumble_id: UUID (FK→Grumble) - 共感対象の投稿ID
- user_id: UUID (FK→AnonymousUser) - 共感した匿名ユーザーID
- vibe_type: VARCHAR(20) (NOT NULL) - 共感の種類（'WAKARU'またはスタンプ名）
- voted_at: TIMESTAMP WITH TIME ZONE (NOT NULL) - 共感した時刻

**インデックス**: grumble_id, (grumble_id, user_id) UNIQUE制約

### 3. 匿名ユーザーテーブル (AnonymousUser)
- user_id: UUID (PK) - 匿名ユーザーの一意識別子（デバイスIDから生成）
- virtue_points: INT (NOT NULL, DEFAULT 0) - 徳ポイント（共感行為で増加）
- created_at: TIMESTAMP WITH TIME ZONE (NOT NULL) - ユーザー作成日時
- profile_title: VARCHAR(50) (NULL) - 称号（例：「今週の菩薩」）

**インデックス**: virtue_points

### 4. 投票テーブル (Poll) - 将来機能
- poll_id: SERIAL (PK) - 投票の一意識別子
- grumble_id: UUID (FK→Grumble) - 投票が紐づく愚痴投稿ID
- question: VARCHAR(255) (NOT NULL) - 投票の質問文
- option_1: VARCHAR(100) (NOT NULL) - 選択肢Aのテキスト
- option_2: VARCHAR(100) (NOT NULL) - 選択肢Bのテキスト

### 5. 投票結果テーブル (PollVote) - 将来機能
- poll_vote_id: SERIAL (PK) - 投票結果の一意識別子
- poll_id: INT (FK→Poll) - 対象の投票ID
- user_id: UUID (FK→AnonymousUser) - 投票した匿名ユーザーID
- selected_option: INT (NOT NULL, 1または2) - 選択した選択肢のインデックス

**インデックス**: (poll_id, user_id) UNIQUE制約

### 6. イベントテーブル (Event)
- event_id: SERIAL (PK) - イベントの一意識別子
- event_name: VARCHAR(100) (NOT NULL) - イベント名（例：「月曜日の大怨霊」）
- event_type: VARCHAR(50) (NOT NULL) - イベント種類（'DAIONRYO'または'OTAKINAGE'）
- start_time: TIMESTAMP WITH TIME ZONE (NOT NULL) - イベント開始時刻
- end_time: TIMESTAMP WITH TIME ZONE (NOT NULL) - イベント終了時刻
- current_hp: INT (NOT NULL) - イベント進行度（HPゲージ、大怨霊のみ）
- max_hp: INT (NOT NULL) - 大怨霊の初期HP
- is_active: BOOLEAN (NOT NULL, DEFAULT FALSE) - 開催中フラグ

## 主要機能
- 24時間自動削除システム
- 匿名ユーザー管理
- 共感（「わかる…」）システム
- 菩薩ランキング（徳ポイント制）
- イベント機能（大怨霊討伐、お焚き上げ）
- 投票機能（将来実装予定）