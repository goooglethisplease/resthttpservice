package http

import (
	"restservice/internal/entity"
)

type subscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type subscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type summaryResponse struct {
	Total int `json:"total"`
}

func toSubscriptionResponse(sub entity.Subscription) subscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		formatted := sub.EndDate.Format(dateLayout)
		endDate = &formatted
	}

	return subscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format(dateLayout),
		EndDate:     endDate,
	}
}
