-- View: grumble_stats_toxic
-- Aggregates per day/week/month per toxic_level for purified/unpurified counts and total vibes.

CREATE OR REPLACE VIEW grumble_stats_toxic AS
SELECT 'day' AS granularity,
       date_trunc('day', posted_at AT TIME ZONE 'UTC') AS bucket,
       toxic_level,
       COUNT(*) FILTER (WHERE is_purified)             AS purified_count,
       COUNT(*) FILTER (WHERE NOT is_purified)         AS unpurified_count,
       SUM(vibe_count)                                 AS total_vibes
  FROM grumbles
 GROUP BY 1,2,3
UNION ALL
SELECT 'week',
       date_trunc('week', posted_at AT TIME ZONE 'UTC'),
       toxic_level,
       COUNT(*) FILTER (WHERE is_purified),
       COUNT(*) FILTER (WHERE NOT is_purified),
       SUM(vibe_count)
  FROM grumbles
 GROUP BY 1,2,3
UNION ALL
SELECT 'month',
       date_trunc('month', posted_at AT TIME ZONE 'UTC'),
       toxic_level,
       COUNT(*) FILTER (WHERE is_purified),
       COUNT(*) FILTER (WHERE NOT is_purified),
       SUM(vibe_count)
  FROM grumbles
 GROUP BY 1,2,3;
