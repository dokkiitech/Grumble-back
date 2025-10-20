-- 初期スキーマ: AnonymousUser, Grumble, Vibe, Poll, PollVote, Event

CREATE TABLE IF NOT EXISTS anonymous_users (
    user_id UUID PRIMARY KEY,
    virtue_points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    profile_title VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS grumbles (
    grumble_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES anonymous_users(user_id),
    content TEXT NOT NULL,
    toxic_level INTEGER NOT NULL CHECK (toxic_level BETWEEN 1 AND 5),
    vibe_count INTEGER NOT NULL DEFAULT 0,
    is_purified BOOLEAN NOT NULL DEFAULT FALSE,
    posted_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_event_grumble BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS vibes (
    vibe_id BIGSERIAL PRIMARY KEY,
    grumble_id UUID NOT NULL REFERENCES grumbles(grumble_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES anonymous_users(user_id),
    vibe_type VARCHAR(20) NOT NULL,
    voted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS vibes_grumble_user_unique ON vibes (grumble_id, user_id);

CREATE TABLE IF NOT EXISTS polls (
    poll_id BIGSERIAL PRIMARY KEY,
    grumble_id UUID NOT NULL REFERENCES grumbles(grumble_id) ON DELETE CASCADE,
    question VARCHAR(255) NOT NULL,
    option_1 VARCHAR(100) NOT NULL,
    option_2 VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS poll_votes (
    poll_vote_id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT NOT NULL REFERENCES polls(poll_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES anonymous_users(user_id),
    selected_option INTEGER NOT NULL CHECK (selected_option IN (1, 2))
);

CREATE UNIQUE INDEX IF NOT EXISTS poll_votes_poll_user_unique ON poll_votes (poll_id, user_id);

CREATE TABLE IF NOT EXISTS events (
    event_id BIGSERIAL PRIMARY KEY,
    event_name VARCHAR(100) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    current_hp INTEGER NOT NULL,
    max_hp INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);
