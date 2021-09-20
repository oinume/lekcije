package ga_measurement

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type gaMeasurementRepository struct {
	client Client
}

func NewGAMeasurementRepository(client Client) *gaMeasurementRepository {
	return &gaMeasurementRepository{client: client}
}

func (r *gaMeasurementRepository) SendEvent(
	ctx context.Context,
	values *model2.GAMeasurementEvent,
	category, action, label string,
	value int64,
	userID uint32,
) error {
	return r.client.SendEvent(ctx, values, category, action, label, value, userID)
}
