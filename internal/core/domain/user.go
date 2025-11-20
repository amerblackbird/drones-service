package domain

type UserType string

const (
	UserTypeEnduser UserType = "enduser"
	UserTypeDrone   UserType = "drone"
	UserTypeAdmin   UserType = "admin"
)

type User struct {
	BaseModel
	Name              string   `json:"name"`
	Email             *string  `json:"email"`
	Phone             *string  `json:"phone"`
	Locale            *string  `json:"locale"`
	Country           *string  `json:"country"`
	DeviceID          *string  `json:"device_id"`
	NotificationToken *string  `json:"notification_token,omitempty"`
	Type              UserType `json:"user_type"` // e.g., "admin", "enduser", "drone"
	Active            bool     `json:"active"`
	Bio               *string  `json:"bio,omitempty"`
	AvatarUrl         *string  `json:"avatar_url,omitempty"`
	DroneId           *string  `json:"drone_id,omitempty"`
}

type UserDTO struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Email             *string  `json:"email"`
	Phone             *string  `json:"phone"`
	Locale            *string  `json:"locale"`
	Country           *string  `json:"country"`
	DeviceID          *string  `json:"device_id"`
	NotificationToken *string  `json:"notification_token,omitempty"`
	Type              UserType `json:"user_type"` // e.g., "admin", "enduser", "drone"
	Active            bool     `json:"active"`
	Bio               *string  `json:"bio,omitempty"`
	AvatarUrl         *string  `json:"avatar_url,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	DeletedAt         *string  `json:"deleted_at,omitempty"`
	DeletedByID       *string  `json:"deleted_by_id,omitempty"`
	CreatedByID       *string  `json:"created_by_id,omitempty"`
	UpdatedByID       *string  `json:"updated_by_id,omitempty"`
}

func (d *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:                d.ID,
		Name:              d.Name,
		Email:             d.Email,
		Phone:             d.Phone,
		Locale:            d.Locale,
		Country:           d.Country,
		DeviceID:          d.DeviceID,
		NotificationToken: d.NotificationToken,
		Type:              d.Type,
		Active:            d.Active,
		Bio:               d.Bio,
		AvatarUrl:         d.AvatarUrl,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		CreatedByID:       d.CreatedByID,
		UpdatedByID:       d.UpdatedByID,
	}
}
