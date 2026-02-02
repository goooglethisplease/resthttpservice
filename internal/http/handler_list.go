package http

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"restservice/internal/usecase/subscription"
)

// ListSubscriptions godoc
// @Summary Список подписок
// @Description Возвращает список подписок с фильтрами.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {array} subscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions [get]
func (h *Handler) handleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("list subscriptions request", "method", r.Method, "path", r.URL.Path)
	query := r.URL.Query()
	userID := strings.TrimSpace(query.Get("user_id"))
	serviceName := strings.TrimSpace(query.Get("service_name"))

	var (
		uid  *uuid.UUID
		name *string
	)

	if userID != "" {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		uid = &parsed
	}

	if serviceName != "" {
		name = &serviceName
	}

	items, err := h.service.List(r.Context(), subscription.ListFilter{
		UserID:      uid,
		ServiceName: name,
	})
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := make([]subscriptionResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toSubscriptionResponse(item))
	}

	h.writeJSON(w, http.StatusOK, resp)
}
