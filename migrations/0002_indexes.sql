-- パフォーマンス向上のためのインデックス

CREATE INDEX IF NOT EXISTS idx_grumbles_posted_at ON grumbles (posted_at);
CREATE INDEX IF NOT EXISTS idx_grumbles_expires_at ON grumbles (expires_at);
CREATE INDEX IF NOT EXISTS idx_grumbles_is_purified ON grumbles (is_purified);

CREATE INDEX IF NOT EXISTS idx_vibes_grumble_id ON vibes (grumble_id);

CREATE INDEX IF NOT EXISTS idx_anonymous_users_virtue_points ON anonymous_users (virtue_points DESC);
