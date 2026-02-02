package http

import (
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"

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

	mux.HandleFunc("/swagger/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Swagger/swagger.yaml")
	})

	mux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Swagger/swagger.yaml")
	})

	mux.Handle(swaggerPath, httpSwagger.WrapHandler)
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
