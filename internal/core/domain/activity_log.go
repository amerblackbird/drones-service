package domain

type LogActivityMetadata struct {
	IP       string `json:"ip"`
	Device   string `json:"device"`
	Location string `json:"location"`
}

type ActivityLog struct {
	BaseModel
	ActorId      string  `json:"actor_id"`
	ActorType    string  `json:"actor_type"`
	Action       string  `json:"action"`
	PerformedAt  string  `json:"performed_at"`
	IP           string  `json:"ip"`
	Device       string  `json:"device"`
	Location     string  `json:"location"`
	ResourceName string  `json:"resource_name"`
	ResourceID   *string `json:"resource_id"`
}

type ActivityLogDTO struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Action       string  `json:"action"`
	PerformedAt  string  `json:"performed_at"`
	ResourceName string  `json:"resource_name"`
	ResourceID   *string `json:"resource_id"`
	IP           string  `json:"ip"`
	Device       string  `json:"device"`
	Location     string  `json:"location"`
	DeletedAt    *string `json:"deleted_at"`
	DeletedByID  *string `json:"deleted_by_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Active       bool    `json:"active"`
	CreatedByID  *string `json:"created_by_id"`
	UpdatedByID  *string `json:"updated_by_id"`
}

type CreateActivityLogRequest struct {
	UserID   string              `json:"user_id"`
	Action   string              `json:"action"`
	Metadata LogActivityMetadata `json:"metadata"`
}

type UpdateActivityLogRequest struct {
	ActorType    *string              `json:"actor_type,omitempty"`
	Action       *string              `json:"action,omitempty"`
	PerformedAt  *string              `json:"performed_at,omitempty"`
	ResourceName *string              `json:"resource_name,omitempty"`
	ResourceID   *string              `json:"resource_id,omitempty"`
	Metadata     *LogActivityMetadata `json:"metadata,omitempty"`
	UpdatedByID  *string              `json:"updated_by_id,omitempty"`
}

type ActivityLogFilter struct {
	UserID    *string `json:"user_id,omitempty"`
	Action    *string `json:"action,omitempty"`
	StartTime *int64  `json:"start_time,omitempty"`
	EndTime   *int64  `json:"end_time,omitempty"`
}

func (a *ActivityLog) ToDTO() *ActivityLogDTO {
	return &ActivityLogDTO{
		ID:           a.ID,
		UserID:       a.ActorId,
		Action:       a.Action,
		PerformedAt:  a.PerformedAt,
		ResourceName: a.ResourceName,
		ResourceID:   a.ResourceID,
		IP:           a.IP,
		Device:       a.Device,
		Location:     a.Location,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		Active:       a.Active,
		CreatedByID:  a.CreatedByID,
		UpdatedByID:  a.UpdatedByID,
	}
}
