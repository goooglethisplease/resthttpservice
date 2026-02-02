package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"restservice/internal/entity"
	"restservice/internal/usecase/subscription"
)

// CreateSubscription godoc
// @Summary Создать подписку
// @Description Создает запись о подписке пользователя.d
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body subscriptionRequest true "Подписка"
// @Success 201 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions [post]
func (h *Handler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	h.logger.Info("create subscription request", "method", r.Method, "path", r.URL.Path)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Info("create subscription decode failed", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	uid, err := uuid.Parse(strings.TrimSpace(req.UserID))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	parsed, err := time.Parse(dateLayout, strings.TrimSpace(req.StartDate))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid start_date")
		return
	}

	startDate := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	var endDate *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsedEnd, err := time.Parse(dateLayout, strings.TrimSpace(*req.EndDate))
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid end_date")
			return
		}
		end := time.Date(parsedEnd.Year(), parsedEnd.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = &end
	}

	sub := entity.Subscription{
		ServiceName: strings.TrimSpace(req.ServiceName),
		Price:       req.Price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	created, err := h.service.Create(r.Context(), sub)
	if err != nil {
		if errors.Is(err, subscription.ErrValidation) {
			h.writeError(w, http.StatusBadRequest, "validation error")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.writeJSON(w, http.StatusCreated, created)
}
