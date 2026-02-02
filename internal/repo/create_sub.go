package repo

import (
	"context"
	"database/sql"
	"fmt"
	"restservice/internal/entity"
)

func (r *SubscriptionRepo) Create(ctx context.Context, s entity.Subscription) (entity.Subscription, error) {
	const q = `
		insert into subscriptions (service_name, price, user_id, start_date, end_date)
		values ($1, $2, $3, $4, $5)
		returning id
		`

	var endDate sql.NullTime
	if s.EndDate != nil {
		endDate = sql.NullTime{Time: s.EndDate.UTC(), Valid: true}
	}

	err := r.db.QueryRowContext(
		ctx,
		q,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate.UTC(),
		endDate,
	).Scan(&s.ID)

	if err != nil {
		return entity.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	return s, nil
}
