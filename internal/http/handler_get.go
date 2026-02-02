package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"restservice/internal/usecase/subscription"
)

// GetSubscription godoc
// @Summary Получить подписку
// @Description Возвращает подписку по идентификатору.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions/{id} [get]
func (h *Handler) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get subscription request", "method", r.Method, "path", r.URL.Path)
	id := strings.TrimPrefix(r.URL.Path, subscriptionsPath+"/")
	if strings.TrimSpace(id) == "" {
		h.writeError(w, http.StatusBadRequest, "missing id")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	sub, err := h.service.Get(r.Context(), uid)
	if err != nil {
		if errors.Is(err, subscription.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.writeJSON(w, http.StatusOK, sub)
}
