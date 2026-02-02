package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"restservice/internal/usecase/subscription"
)

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `delete from subscriptions where id = $1`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete subscription rows: %w", err)
	}

	if affected == 0 {
		return subscription.ErrNotFound
	}
	return nil
}
