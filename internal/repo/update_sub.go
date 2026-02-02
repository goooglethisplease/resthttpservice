package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"restservice/internal/entity"
	"restservice/internal/usecase/subscription"
)

func (r *SubscriptionRepo) Update(ctx context.Context, id uuid.UUID, s entity.Subscription) (entity.Subscription, error) {
	const q = `
		update subscriptions
		set service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5
		where id = $6
		returning id, service_name, price, user_id, start_date, end_date
	`

	var endDate sql.NullTime
	if s.EndDate != nil {
		endDate = sql.NullTime{Time: s.EndDate.UTC(), Valid: true}
	}

	var (
		sub entity.Subscription
		ed  sql.NullTime
	)

	err := r.db.QueryRowContext(
		ctx,
		q,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate.UTC(),
		endDate,
		id,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&ed,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Subscription{}, subscription.ErrNotFound
		}
		return entity.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}

	sub.StartDate = sub.StartDate.UTC()
	if ed.Valid {
		t := ed.Time.UTC()
		sub.EndDate = &t
	}

	return sub, nil
}
