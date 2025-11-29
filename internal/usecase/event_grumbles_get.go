package usecase

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	sharedservice "github.com/dokkiitech/grumble-back/internal/domain/shared/service"
)

// EventGrumblesGetUseCase handles retrieving event grumbles from archive
type EventGrumblesGetUseCase struct {
	grumbleRepo  grumble.Repository
	eventTimeSvc *sharedservice.EventTimeService
}

// NewEventGrumblesGetUseCase creates a new EventGrumblesGetUseCase
func NewEventGrumblesGetUseCase(grumbleRepo grumble.Repository, eventTimeSvc *sharedservice.EventTimeService) *EventGrumblesGetUseCase {
	return &EventGrumblesGetUseCase{
		grumbleRepo:  grumbleRepo,
		eventTimeSvc: eventTimeSvc,
	}
}

// EventGrumblesRequest represents request parameters
type EventGrumblesRequest struct {
	ToxicLevelMin *shared.ToxicLevel
	ToxicLevelMax *shared.ToxicLevel
	IsPurified    *bool // nilの場合は全て、trueの場合は成仏済み、falseの場合は成仏していない
	Limit         int
	Offset        int
}

// EventGrumblesResponse represents the response
type EventGrumblesResponse struct {
	Grumbles      []*grumble.Grumble
	Total         int
	IsEventActive bool
	EventDate     time.Time
}

// Get retrieves event grumbles if within event time window (00:00-12:00 JST)
func (uc *EventGrumblesGetUseCase) Get(ctx context.Context, req EventGrumblesRequest) (*EventGrumblesResponse, error) {
	now := time.Now()

	// イベント期間中かチェック（24:00〜12:00）
	if !uc.eventTimeSvc.IsEventTimeWindow(now) {
		// イベント期間外の場合、空のレスポンスを返す
		return &EventGrumblesResponse{
			Grumbles:      []*grumble.Grumble{},
			Total:         0,
			IsEventActive: false,
			EventDate:     time.Time{},
		}, nil
	}

	// イベント対象日を取得（前日）
	targetDate := uc.eventTimeSvc.GetEventTargetDate(now)

	// フィルタ構築
	filter := grumble.TimelineFilter{
		ToxicLevelMin:  req.ToxicLevelMin,
		ToxicLevelMax:  req.ToxicLevelMax,
		IsPurified:     req.IsPurified,
		ExcludeExpired: false, // アーカイブなので期限チェック不要
		Limit:          req.Limit,
		Offset:         req.Offset,
	}

	// アーカイブテーブルから前日の投稿を取得
	grumbles, err := uc.grumbleRepo.FindArchivedTimeline(ctx, filter, targetDate)
	if err != nil {
		return nil, err
	}

	total, err := uc.grumbleRepo.CountArchivedTimeline(ctx, filter, targetDate)
	if err != nil {
		return nil, err
	}

	return &EventGrumblesResponse{
		Grumbles:      grumbles,
		Total:         total,
		IsEventActive: true,
		EventDate:     targetDate,
	}, nil
}
