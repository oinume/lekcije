package ga_measurement

import (
	"context"

	"github.com/oinume/lekcije/backend/ga_measurement"
)

type gaMeasurementRepository struct {
	client ga_measurement.Client
}

func NewGAMeasurementRepository(client ga_measurement.Client) *gaMeasurementRepository {
	return &gaMeasurementRepository{client: client}
}

func (r *gaMeasurementRepository) SendEvent(
	ctx context.Context,
	values *ga_measurement.EventValues,
	category, action, label string,
	value int64,
	userID uint32,
) error {
	return r.client.SendEvent(values, category, action, label, value, userID)
}
