package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (r *notificationTimeSpanRepository) UpdateAll(ctx context.Context, userID uint, timeSpans []*model2.NotificationTimeSpan) error {
	for _, ts := range timeSpans {
		if userID != ts.UserID {
			return errors.NewInvalidArgumentError(
				errors.WithMessage("Given userID and userID of timeSpans must be same"),
			)
		}
	}
	if err := repository.Transaction(ctx, r.db, func(exec repository.Executor) error {
		if _, err := model2.NotificationTimeSpans(qm.Where("user_id = ?", userID)).DeleteAll(ctx, exec); err != nil {
			return errors.NewInternalError(
				errors.WithError(err),
				errors.WithMessage("UpdateAll delete failed"),
				errors.WithResource(errors.NewResource("notification_time_spans", "userID", userID)),
			)
		}
		for _, ts := range timeSpans {
			if err := ts.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.NewInternalError(
					errors.WithError(err),
					errors.WithMessage("UpdateAll insert failed"),
					errors.WithResource(errors.NewResource("notification_time_spans", "userID", userID)),
				)
			}
		}
		return nil
	}); err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("UpdateAll commit failed"),
			errors.WithResource(errors.NewResource("notification_time_spans", "userID", userID)),
		)
	}
	return nil
}
