package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"restservice/internal/usecase/subscription"
)

// SummarySubscriptions godoc
// @Summary Сумма подписок
// @Description Считает стоимость подписок за период с фильтрацией.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Дата начала (MM-YYYY)"
// @Param end_date query string true "Дата окончания (MM-YYYY)"
// @Success 200 {object} summaryResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/subscriptions/summary [get]
func (h *Handler) handleSummary(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("summary subscriptions request", "method", r.Method, "path", r.URL.Path)
	query := r.URL.Query()
	userID := strings.TrimSpace(query.Get("user_id"))
	serviceName := strings.TrimSpace(query.Get("service_name"))
	startDateRaw := strings.TrimSpace(query.Get("start_date"))
	endDateRaw := strings.TrimSpace(query.Get("end_date"))

	if startDateRaw == "" || endDateRaw == "" {
		h.writeError(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	startParsed, err := time.Parse(dateLayout, startDateRaw)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid start_date")
		return
	}

	endParsed, err := time.Parse(dateLayout, endDateRaw)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid end_date")
		return
	}

	startDate := time.Date(startParsed.Year(), startParsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(endParsed.Year(), endParsed.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	if endDate.Before(startDate) {
		h.writeError(w, http.StatusBadRequest, "end_date before start_date")
		return
	}

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

	total, err := h.service.Sum(r.Context(), subscription.SummaryFilter{
		UserID:      uid,
		ServiceName: name,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.writeJSON(w, http.StatusOK, summaryResponse{Total: total})
}
