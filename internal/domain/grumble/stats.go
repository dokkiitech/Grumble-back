package grumble

import "time"

// Granularity represents the time bucket size for statistics.
type Granularity string

const (
	GranularityDay   Granularity = "day"
	GranularityWeek  Granularity = "week"
	GranularityMonth Granularity = "month"
)

// StatsRow holds aggregated counts for a time bucket.
type StatsRow struct {
	Bucket          time.Time
	PurifiedCount   int
	UnpurifiedCount int
	TotalVibes      int
	ToxicLevel      *int // Optional; present for toxic-level breakdown
}
