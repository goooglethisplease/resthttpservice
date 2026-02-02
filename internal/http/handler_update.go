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

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет подписку по идентификатору.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body subscriptionRequest true "Подписка"
// @Success 200 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions/{id} [put]
func (h *Handler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("update subscription request", "method", r.Method, "path", r.URL.Path)
	id := strings.TrimPrefix(r.URL.Path, subscriptionsPath+"/")
	if strings.TrimSpace(id) == "" {
		h.writeError(w, http.StatusBadRequest, "missing id")
		return
	}

	subID, err := uuid.Parse(id)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req subscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Info("update subscription decode failed", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	userID, err := uuid.Parse(strings.TrimSpace(req.UserID))
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
		ID:          subID,
		ServiceName: strings.TrimSpace(req.ServiceName),
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	updated, err := h.service.Update(r.Context(), subID, sub)
	if err != nil {
		if errors.Is(err, subscription.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, subscription.ErrValidation) {
			h.writeError(w, http.StatusBadRequest, "validation error")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.writeJSON(w, http.StatusOK, toSubscriptionResponse(updated))
}
