package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"restservice/internal/usecase/subscription"
)

const (
	subscriptionsPath = "/api/subscriptions"
	summaryPath       = "/api/subscriptions/summary"
	swaggerPath       = "/swagger/"
	dateLayout        = "01-2006"
)

type Handler struct {
	service *subscription.Service
	logger  *slog.Logger
}

func NewHandler(service *subscription.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc(summaryPath, h.handleSummary)
	mux.HandleFunc(subscriptionsPath, h.handleSubscriptions)
	mux.HandleFunc(subscriptionsPath+"/", h.handleSubscriptionByID)
	mux.Handle(swaggerPath, httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("encode response failed", "error", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, errorResponse{Error: message})
}
