package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"restservice/internal/usecase/subscription"
)

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по идентификатору.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions/{id} [delete]
func (h *Handler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("delete subscription request", "method", r.Method, "path", r.URL.Path)
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

	if err := h.service.Delete(r.Context(), subID); err != nil {
		if errors.Is(err, subscription.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
