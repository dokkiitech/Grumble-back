package service

import "time"

// EventTimeService handles event time window logic and time calculations.
type EventTimeService struct {
	timezone       *time.Location
	eventStartHour int // イベント開始時刻（00:00）
	eventEndHour   int // イベント終了時刻（12:00）
}

// NewEventTimeService creates a new EventTimeService with JST timezone.
func NewEventTimeService() *EventTimeService {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// フォールバック: UTCで計算
		jst = time.FixedZone("JST", 9*60*60)
	}

	return &EventTimeService{
		timezone:       jst,
		eventStartHour: 0,  // 00:00
		eventEndHour:   12, // 12:00
	}
}

// IsEventTimeWindow checks if current time is in event display window (00:00-12:00 JST).
func (s *EventTimeService) IsEventTimeWindow(now time.Time) bool {
	nowJST := now.In(s.timezone)
	hour := nowJST.Hour()

	// 00:00 〜 11:59 の間
	return hour >= s.eventStartHour && hour < s.eventEndHour
}

// GetEventTargetDate returns the date for which archived grumbles should be displayed.
// If current time is 00:00-11:59, returns previous day's date.
func (s *EventTimeService) GetEventTargetDate(now time.Time) time.Time {
	nowJST := now.In(s.timezone)

	// 00:00-11:59の場合、前日の投稿を表示
	if nowJST.Hour() < s.eventEndHour {
		return nowJST.AddDate(0, 0, -1)
	}

	// 12:00以降は当日（ただしイベント期間外なので使われない）
	return nowJST
}

// GetDayBounds returns the start and end time of a given date in JST.
func (s *EventTimeService) GetDayBounds(date time.Time) (start, end time.Time) {
	dateJST := date.In(s.timezone)

	// その日の00:00:00
	start = time.Date(
		dateJST.Year(),
		dateJST.Month(),
		dateJST.Day(),
		0, 0, 0, 0,
		s.timezone,
	)

	// その日の23:59:59
	end = start.Add(24*time.Hour - time.Second)

	return start, end
}

// CalculateNextMidnight calculates the next midnight (00:00) in JST.
// Used for setting expiration time of grumbles.
func (s *EventTimeService) CalculateNextMidnight(now time.Time) time.Time {
	nowJST := now.In(s.timezone)

	// 当日の00:00
	midnight := time.Date(
		nowJST.Year(),
		nowJST.Month(),
		nowJST.Day(),
		0, 0, 0, 0,
		s.timezone,
	)

	// 翌日の00:00
	return midnight.Add(24 * time.Hour)
}
