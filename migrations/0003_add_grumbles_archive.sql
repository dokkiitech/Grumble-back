-- アーカイブテーブル（イベント用）
-- 24:00を過ぎた投稿をgrumblesテーブルから移動して保存

CREATE TABLE IF NOT EXISTS grumbles_archive (
    grumble_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    content TEXT NOT NULL CHECK (length(content) <= 280),
    toxic_level INTEGER NOT NULL CHECK (toxic_level BETWEEN 1 AND 5),
    vibe_count INTEGER NOT NULL DEFAULT 0,
    is_purified BOOLEAN NOT NULL DEFAULT FALSE,
    posted_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_event_grumble BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- インデックス
CREATE INDEX idx_grumbles_archive_posted_at ON grumbles_archive(posted_at DESC);
CREATE INDEX idx_grumbles_archive_archived_at ON grumbles_archive(archived_at DESC);
CREATE INDEX idx_grumbles_archive_user_id ON grumbles_archive(user_id);
CREATE INDEX idx_grumbles_archive_is_event_grumble ON grumbles_archive(is_event_grumble) WHERE is_event_grumble = TRUE;
CREATE INDEX idx_grumbles_archive_toxic_level ON grumbles_archive(toxic_level);
