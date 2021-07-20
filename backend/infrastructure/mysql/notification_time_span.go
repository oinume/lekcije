package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opencensus.io/trace"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type notificationTimeSpanRepository struct {
	db *sql.DB
}

func NewNotificationTimeSpanRepository(db *sql.DB) repository.NotificationTimeSpan {
	return &notificationTimeSpanRepository{db: db}
}

func (r *notificationTimeSpanRepository) FindByUserID(ctx context.Context, userID uint) ([]*model2.NotificationTimeSpan, error) {
	_, span := trace.StartSpan(ctx, "NotificationTimeSpanService.FindByUserID")
	defer span.End()
	span.Annotatef([]trace.Attribute{
		trace.Int64Attribute("userID", int64(userID)),
	}, "userID:%d", userID)

	timeSpans, err := model2.NotificationTimeSpans(qm.Where("user_id = ?", userID)).All(ctx, r.db)
	if err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("FindByUserID select failed"),
			errors.WithResource(errors.NewResource((&model.NotificationTimeSpan{}).TableName(), "userID", userID)),
		)
	}
	return timeSpans, nil
}
