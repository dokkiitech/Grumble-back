package service

import (
	"testing"
	"time"
)

func TestEventTimeService_IsEventTimeWindow(t *testing.T) {
	svc := NewEventTimeService()

	tests := []struct {
		name     string
		hour     int
		minute   int
		expected bool
	}{
		{"00:00 - イベント期間内", 0, 0, true},
		{"01:30 - イベント期間内", 1, 30, true},
		{"11:59 - イベント期間内", 11, 59, true},
		{"12:00 - イベント期間外", 12, 0, false},
		{"13:00 - イベント期間外", 13, 0, false},
		{"23:59 - イベント期間外", 23, 59, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// JSTで特定の時刻を作成
			jst, _ := time.LoadLocation("Asia/Tokyo")
			testTime := time.Date(2025, 11, 27, tt.hour, tt.minute, 0, 0, jst)

			result := svc.IsEventTimeWindow(testTime)
			if result != tt.expected {
				t.Errorf("IsEventTimeWindow(%02d:%02d) = %v, want %v",
					tt.hour, tt.minute, result, tt.expected)
			}
		})
	}
}

func TestEventTimeService_GetEventTargetDate(t *testing.T) {
	svc := NewEventTimeService()
	jst, _ := time.LoadLocation("Asia/Tokyo")

	tests := []struct {
		name           string
		currentTime    time.Time
		expectedDayDiff int // 現在日との差分（-1なら前日）
	}{
		{
			name:           "午前0時 - 前日を返す",
			currentTime:    time.Date(2025, 11, 27, 0, 0, 0, 0, jst),
			expectedDayDiff: -1,
		},
		{
			name:           "午前11時59分 - 前日を返す",
			currentTime:    time.Date(2025, 11, 27, 11, 59, 0, 0, jst),
			expectedDayDiff: -1,
		},
		{
			name:           "午後12時 - 当日を返す",
			currentTime:    time.Date(2025, 11, 27, 12, 0, 0, 0, jst),
			expectedDayDiff: 0,
		},
		{
			name:           "午後23時 - 当日を返す",
			currentTime:    time.Date(2025, 11, 27, 23, 0, 0, 0, jst),
			expectedDayDiff: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDate := svc.GetEventTargetDate(tt.currentTime)
			expected := tt.currentTime.AddDate(0, 0, tt.expectedDayDiff)

			// 日付のみ比較
			if targetDate.Year() != expected.Year() ||
				targetDate.Month() != expected.Month() ||
				targetDate.Day() != expected.Day() {
				t.Errorf("GetEventTargetDate(%v) = %v, want %v",
					tt.currentTime.Format("2006-01-02 15:04"),
					targetDate.Format("2006-01-02"),
					expected.Format("2006-01-02"))
			}
		})
	}
}

func TestEventTimeService_GetDayBounds(t *testing.T) {
	svc := NewEventTimeService()
	jst, _ := time.LoadLocation("Asia/Tokyo")

	testDate := time.Date(2025, 11, 27, 15, 30, 0, 0, jst)
	start, end := svc.GetDayBounds(testDate)

	// 開始時刻: 2025-11-27 00:00:00
	expectedStart := time.Date(2025, 11, 27, 0, 0, 0, 0, jst)
	if !start.Equal(expectedStart) {
		t.Errorf("GetDayBounds start = %v, want %v", start, expectedStart)
	}

	// 終了時刻: 2025-11-27 23:59:59
	expectedEnd := time.Date(2025, 11, 27, 23, 59, 59, 0, jst)
	if !end.Equal(expectedEnd) {
		t.Errorf("GetDayBounds end = %v, want %v", end, expectedEnd)
	}
}

func TestEventTimeService_CalculateNextMidnight(t *testing.T) {
	svc := NewEventTimeService()
	jst, _ := time.LoadLocation("Asia/Tokyo")

	tests := []struct {
		name        string
		currentTime time.Time
		expected    time.Time
	}{
		{
			name:        "午前0時 - 翌日0時を返す",
			currentTime: time.Date(2025, 11, 27, 0, 0, 0, 0, jst),
			expected:    time.Date(2025, 11, 28, 0, 0, 0, 0, jst),
		},
		{
			name:        "午後15時30分 - 翌日0時を返す",
			currentTime: time.Date(2025, 11, 27, 15, 30, 0, 0, jst),
			expected:    time.Date(2025, 11, 28, 0, 0, 0, 0, jst),
		},
		{
			name:        "午後23時59分 - 翌日0時を返す",
			currentTime: time.Date(2025, 11, 27, 23, 59, 0, 0, jst),
			expected:    time.Date(2025, 11, 28, 0, 0, 0, 0, jst),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.CalculateNextMidnight(tt.currentTime)
			if !result.Equal(tt.expected) {
				t.Errorf("CalculateNextMidnight(%v) = %v, want %v",
					tt.currentTime.Format("2006-01-02 15:04"),
					result.Format("2006-01-02 15:04"),
					tt.expected.Format("2006-01-02 15:04"))
			}
		})
	}
}
