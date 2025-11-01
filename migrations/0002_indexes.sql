-- Grumbles indexes for performance
CREATE INDEX IF NOT EXISTS idx_grumbles_posted_at ON grumbles(posted_at DESC);
CREATE INDEX IF NOT EXISTS idx_grumbles_expires_at ON grumbles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_grumbles_is_purified ON grumbles(is_purified) WHERE is_purified = FALSE;
CREATE INDEX IF NOT EXISTS idx_grumbles_toxic_level ON grumbles(toxic_level);

-- Vibes indexes
CREATE INDEX IF NOT EXISTS idx_vibes_grumble_id ON vibes(grumble_id);

-- Anonymous users indexes
CREATE INDEX IF NOT EXISTS idx_anonymous_users_virtue_points ON anonymous_users(virtue_points DESC);

-- Events indexes
CREATE INDEX IF NOT EXISTS idx_events_is_active ON events(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_events_end_time ON events(end_time);
