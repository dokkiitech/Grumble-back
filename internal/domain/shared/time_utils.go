package shared

import "time"

// IsEventTimeWindow checks if current time is in event display window (00:00-12:00 JST)
func IsEventTimeWindow(now time.Time) bool {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// フォールバック: UTCで計算
		jst = time.FixedZone("JST", 9*60*60)
	}

	nowJST := now.In(jst)
	hour := nowJST.Hour()

	// 00:00 〜 11:59 の間
	return hour >= 0 && hour < 12
}

// GetEventTargetDate returns the date for which archived grumbles should be displayed
// If current time is 00:00-11:59, returns previous day's date
func GetEventTargetDate(now time.Time) time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", 9*60*60)
	}

	nowJST := now.In(jst)

	// 00:00-11:59の場合、前日の投稿を表示
	if nowJST.Hour() < 12 {
		return nowJST.AddDate(0, 0, -1)
	}

	// 12:00以降は当日（ただしイベント期間外なので使われない）
	return nowJST
}

// GetDayBounds returns the start and end time of a given date in JST
func GetDayBounds(date time.Time) (start, end time.Time) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", 9*60*60)
	}

	dateJST := date.In(jst)

	// その日の00:00:00
	start = time.Date(
		dateJST.Year(),
		dateJST.Month(),
		dateJST.Day(),
		0, 0, 0, 0,
		jst,
	)

	// その日の23:59:59
	end = start.Add(24*time.Hour - time.Second)

	return start, end
}

// CalculateNextMidnight calculates the next midnight (00:00) in JST
// Used for setting expiration time of grumbles
func CalculateNextMidnight(now time.Time) time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", 9*60*60)
	}

	nowJST := now.In(jst)

	// 当日の00:00
	midnight := time.Date(
		nowJST.Year(),
		nowJST.Month(),
		nowJST.Day(),
		0, 0, 0, 0,
		jst,
	)

	// 翌日の00:00
	return midnight.Add(24 * time.Hour)
}
