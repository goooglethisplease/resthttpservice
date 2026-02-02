package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"restservice/internal/entity"
	"restservice/internal/usecase/subscription"
)

type memoryRepo struct {
	items map[uuid.UUID]entity.Subscription
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{items: make(map[uuid.UUID]entity.Subscription)}
}

func (r *memoryRepo) Create(ctx context.Context, s entity.Subscription) (entity.Subscription, error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	r.items[s.ID] = s
	return s, nil
}

func (r *memoryRepo) Get(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {
	item, ok := r.items[id]
	if !ok {
		return entity.Subscription{}, subscription.ErrNotFound
	}
	return item, nil
}

func (r *memoryRepo) List(ctx context.Context, filter subscription.ListFilter) ([]entity.Subscription, error) {
	result := make([]entity.Subscription, 0)
	for _, item := range r.items {
		if filter.UserID != nil && item.UserID != *filter.UserID {
			continue
		}
		if filter.ServiceName != nil && item.ServiceName != *filter.ServiceName {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *memoryRepo) Update(ctx context.Context, id uuid.UUID, s entity.Subscription) (entity.Subscription, error) {
	if _, ok := r.items[id]; !ok {
		return entity.Subscription{}, subscription.ErrNotFound
	}
	s.ID = id
	r.items[id] = s
	return s, nil
}

func (r *memoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := r.items[id]; !ok {
		return subscription.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func (r *memoryRepo) Sum(ctx context.Context, filter subscription.SummaryFilter) (int, error) {
	total := 0
	for _, item := range r.items {
		if filter.UserID != nil && item.UserID != *filter.UserID {
			continue
		}
		if filter.ServiceName != nil && item.ServiceName != *filter.ServiceName {
			continue
		}
		if item.StartDate.After(filter.EndDate) {
			continue
		}
		if item.EndDate != nil && item.EndDate.Before(filter.StartDate) {
			continue
		}
		total += item.Price
	}
	return total, nil
}

func newTestHandler(repo *memoryRepo) *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
	service := subscription.NewService(repo, logger)
	return NewHandler(service, logger)
}

func TestCreateSubscription(t *testing.T) {
	repo := newMemoryRepo()
	handler := newTestHandler(repo)

	body := subscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, subscriptionsPath, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	handler.handleCreateSubscription(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	var resp subscriptionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.ID == "" {
		t.Fatal("expected id")
	}
	if resp.ServiceName != body.ServiceName {
		t.Fatalf("unexpected service name: %s", resp.ServiceName)
	}
	if resp.Price != body.Price {
		t.Fatalf("unexpected price: %d", resp.Price)
	}
	if resp.UserID != body.UserID {
		t.Fatalf("unexpected user id: %s", resp.UserID)
	}
	if resp.StartDate != body.StartDate {
		t.Fatalf("unexpected start date: %s", resp.StartDate)
	}
}

func TestSummarySubscriptions(t *testing.T) {
	repo := newMemoryRepo()
	handler := newTestHandler(repo)

	userID := uuid.New()
	repo.items[uuid.New()] = entity.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix",
		Price:       500,
		UserID:      userID,
		StartDate:   time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.items[uuid.New()] = entity.Subscription{
		ID:          uuid.New(),
		ServiceName: "Spotify",
		Price:       300,
		UserID:      userID,
		StartDate:   time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     func() *time.Time { t := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC); return &t }(),
	}

	req := httptest.NewRequest(http.MethodGet, summaryPath+"?start_date=07-2025&end_date=08-2025&user_id="+userID.String(), nil)
	w := httptest.NewRecorder()

	handler.handleSummary(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	var resp summaryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Total != 800 {
		t.Fatalf("unexpected total: %d", resp.Total)
	}
}
