package subscription

import (
	"time"

	"github.com/google/uuid"
)

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
}

type SummaryFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   time.Time
	EndDate     time.Time
}
