package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"restservice/internal/entity"
)

const (
	minPrice = 1
)

var (
	ErrNotFound   = errors.New("subscription not found")
	ErrValidation = errors.New("validation error")
)

type Repository interface {
	Create(ctx context.Context, s entity.Subscription) (entity.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Subscription, error)
	List(ctx context.Context, filter ListFilter) ([]entity.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, s entity.Subscription) (entity.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Sum(ctx context.Context, filter SummaryFilter) (int, error)
}

type Service struct {
	repository Repository
	logger     *slog.Logger
}

func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}

func Validate(s entity.Subscription) error {
	if strings.TrimSpace(s.ServiceName) == "" {
		return errors.New("service name is required")
	}
	if s.Price < minPrice {
		return errors.New("price too low")
	}
	if s.UserID == uuid.Nil {
		return errors.New("user id is required")
	}
	if s.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return errors.New("end date be before start date")
	}
	return nil
}

func (s *Service) Create(ctx context.Context, sub entity.Subscription) (entity.Subscription, error) {
	if err := Validate(sub); err != nil {
		s.logger.Info("subscription validation failed", "error", err)
		return entity.Subscription{}, fmt.Errorf("%w: %s", ErrValidation, err)
	}

	created, err := s.repository.Create(ctx, sub)
	if err != nil {
		s.logger.Error("subscription create failed", "error", err)
		return entity.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	s.logger.Info("subscription created", "subscription_id", created.ID)

	return created, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {
	sub, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Info("subscription not found", "subscription_id", id)
			return entity.Subscription{}, ErrNotFound
		}
		s.logger.Error("subscription get failed", "error", err, "subscription_id", id)
		return entity.Subscription{}, fmt.Errorf("get subscription: %w", err)
	}

	s.logger.Info("subscription fetched", "subscription_id", id)
	return sub, nil
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]entity.Subscription, error) {
	items, err := s.repository.List(ctx, filter)
	if err != nil {
		s.logger.Error("subscription list failed", "error", err)
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	s.logger.Info("subscription list fetched", "count", len(items))
	return items, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, sub entity.Subscription) (entity.Subscription, error) {
	if err := Validate(sub); err != nil {
		s.logger.Info("subscription validation failed", "error", err)
		return entity.Subscription{}, fmt.Errorf("%w: %s", ErrValidation, err)
	}

	updated, err := s.repository.Update(ctx, id, sub)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Info("subscription not found", "subscription_id", id)
			return entity.Subscription{}, ErrNotFound
		}
		s.logger.Error("subscription update failed", "error", err, "subscription_id", id)
		return entity.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}

	s.logger.Info("subscription updated", "subscription_id", id)
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Info("subscription not found", "subscription_id", id)
			return ErrNotFound
		}
		s.logger.Error("subscription delete failed", "error", err, "subscription_id", id)
		return fmt.Errorf("delete subscription: %w", err)
	}

	s.logger.Info("subscription deleted", "subscription_id", id)
	return nil
}

func (s *Service) Sum(ctx context.Context, filter SummaryFilter) (int, error) {
	total, err := s.repository.Sum(ctx, filter)
	if err != nil {
		s.logger.Error("subscription summary failed", "error", err)
		return 0, fmt.Errorf("summary subscriptions: %w", err)
	}

	s.logger.Info("subscription summary fetched", "total", total)
	return total, nil
}
