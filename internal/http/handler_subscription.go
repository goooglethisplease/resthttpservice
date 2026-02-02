package http

import "net/http"

func (h *Handler) handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateSubscription(w, r)
	case http.MethodGet:
		h.handleListSubscriptions(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetSubscription(w, r)
	case http.MethodPut:
		h.handleUpdateSubscription(w, r)
	case http.MethodDelete:
		h.handleDeleteSubscription(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
