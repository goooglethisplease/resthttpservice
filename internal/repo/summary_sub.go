package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"restservice/internal/usecase/subscription"
)

func (r *SubscriptionRepo) Sum(ctx context.Context, filter subscription.SummaryFilter) (int, error) {
	base := `
		select coalesce(sum(price), 0)
		from subscriptions
		where start_date <= $1
		  and (end_date is null or end_date >= $2)
	`

	args := []any{filter.EndDate.UTC(), filter.StartDate.UTC()}
	cond := make([]string, 0, 2)

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

	var total int
	if err := r.db.QueryRowContext(ctx, base, args...).Scan(&total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("sum subscriptions: %w", err)
	}

	return total, nil
}
