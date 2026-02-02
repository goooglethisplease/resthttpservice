package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"restservice/internal/entity"
	"restservice/internal/usecase/subscription"
)

func (r *SubscriptionRepo) List(ctx context.Context, filter subscription.ListFilter) ([]entity.Subscription, error) {
	base := `
		select id, service_name, price, user_id, start_date, end_date
		from subscriptions
		where 1=1
	`
	var (
		args []any
		cond []string
	)

	if filter.UserID != nil {
		args = append(args, *filter.UserID)
		cond = append(cond, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if filter.ServiceName != nil && strings.TrimSpace(*filter.ServiceName) != "" {
		args = append(args, strings.TrimSpace(*filter.ServiceName))
		cond = append(cond, fmt.Sprintf("service_name = $%d", len(args)))
	}

	if len(cond) > 0 {
		base += " and " + strings.Join(cond, " and ")
	}
	base += " order by start_date desc, id"

	rows, err := r.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entity.Subscription
	for rows.Next() {
		var (
			sub     entity.Subscription
			endDate sql.NullTime
		)
		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDate,
		); err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		sub.StartDate = sub.StartDate.UTC()
		if endDate.Valid {
			t := endDate.Time.UTC()
			sub.EndDate = &t
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list subscriptions rows: %w", err)
	}

	return subs, nil
}
