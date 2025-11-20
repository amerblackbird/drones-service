package domain

type AuditLog struct {
	BaseModel
	PerformedById     string      `json:"performed_by_id"`
	PerformedUserType string      `json:"performed_user_type"`
	Action            string      `json:"action"`
	ResourceName      string      `json:"resource_name"`
	ResourceID        string      `json:"resource_id"`
	PerformedAt       string      `json:"performed_at"`
	OldData           interface{} `json:"old_data"`
	NewData           interface{} `json:"new_data"`
}

type AuditLogDTO struct {
	ID                string      `json:"id"`
	PerformedById     string      `json:"performed_by_id"`
	PerformedUserType string      `json:"performed_user_type"`
	Action            string      `json:"action"`
	ResourceName      string      `json:"resource_name"`
	ResourceID        string      `json:"resource_id"`
	PerformedAt       string      `json:"performed_at"`
	OldData           interface{} `json:"old_data"`
	NewData           interface{} `json:"new_data"`
	DeletedAt         *string     `json:"deleted_at"`
	DeletedByID       *string     `json:"deleted_by_id"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
	Active            bool        `json:"active"`
	CreatedByID       *string     `json:"created_by_id"`
	UpdatedByID       *string     `json:"updated_by_id"`
}

type CreateAuditLogRequest struct {
	PerformedBy  string      `json:"performed_by_id"`
	UserType     string      `json:"user_type"`
	Action       string      `json:"action"`
	ResourceName string      `json:"resource_name"`
	ResourceID   string      `json:"resource_id"`
	PerformedAt  string      `json:"performed_at"`
	OldData      interface{} `json:"old_data"`
	NewData      interface{} `json:"new_data"`
}

type UpdateAuditLogRequest struct {
	PerformedUserType *string     `json:"performed_user_type,omitempty"`
	Action            *string     `json:"action,omitempty"`
	ResourceName      *string     `json:"resource_name,omitempty"`
	ResourceID        *string     `json:"resource_id,omitempty"`
	PerformedAt       *string     `json:"performed_at,omitempty"`
	OldData           interface{} `json:"old_data,omitempty"`
	NewData           interface{} `json:"new_data,omitempty"`
	UpdatedByID       *string     `json:"updated_by_id,omitempty"`
}

type AuditLogFilter struct {
	PerformedBy  *string `json:"performed_by_id,omitempty"`
	UserType     *string `json:"user_type,omitempty"`
	Action       *string `json:"action,omitempty"`
	ResourceName *string `json:"resource_name,omitempty"`
	ResourceID   *string `json:"resource_id,omitempty"`
}

func (a *AuditLog) ToDTO() *AuditLogDTO {
	return &AuditLogDTO{
		ID:                a.ID,
		PerformedById:     a.PerformedById,
		PerformedUserType: a.PerformedUserType,
		Action:            a.Action,
		ResourceName:      a.ResourceName,
		ResourceID:        a.ResourceID,
		PerformedAt:       a.PerformedAt,
		OldData:           a.OldData,
		NewData:           a.NewData,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		Active:            a.Active,
		CreatedByID:       a.CreatedByID,
		UpdatedByID:       a.UpdatedByID,
	}
}
