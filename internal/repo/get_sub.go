package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"restservice/internal/entity"
	"restservice/internal/usecase/subscription"
)

func (r *SubscriptionRepo) Get(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {
	const q = `
		select id, service_name, price, user_id, start_date, end_date
		from subscriptions
		where id = $1
	`

	var (
		sub     entity.Subscription
		endDate sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&endDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Subscription{}, subscription.ErrNotFound
		}
		return entity.Subscription{}, err
	}

	sub.StartDate = sub.StartDate.UTC()

	if endDate.Valid {
		t := endDate.Time.UTC()
		sub.EndDate = &t
	}

	return sub, nil
}
