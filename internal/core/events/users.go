package events

import "drones/internal/core/domain"

type UserCreatedEvent struct {
	UserID            string          `json:"user_id"`
	Name              string          `json:"name"`
	Email             *string         `json:"email"`
	Phone             *string         `json:"phone"`
	Locale            *string         `json:"locale"`
	Country           *string         `json:"country"`
	DeviceID          *string         `json:"device_id"`
	NotificationToken *string         `json:"notification_token,omitempty"`
	UserType          domain.UserType `json:"user_type"` // e.g., "customer", "driver"
}

type UserUpdatedEvent struct {
	UserID    string  `json:"user_id"`
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Locale    *string `json:"locale"`
	Country   *string `json:"country"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
}

type UserLoggedInEvent struct {
	UserID            string  `json:"user_id"`
	DeviceID          *string `json:"device_id"`
	NotificationToken *string `json:"notification_token,omitempty"`
	BrowserInfo       *string `json:"browser_info,omitempty"`
	IPAddress         *string `json:"ip_address,omitempty"`
	Location          *string `json:"location,omitempty"`
}
