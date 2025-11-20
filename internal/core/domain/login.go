package domain

type Login struct {
	ID                 string   `json:"login_id"`
	UserID             string   `json:"user_id"`
	UserType           UserType `json:"user_type,omitempty"`
	Username           string   `json:"username"`
	IPAddress          *string  `json:"ip_address"`
	UserAgent          *string  `json:"user_agent,omitempty"`
	HashedPassword     *string  `json:"hashed_password"`
	Salt               *string  `json:"salt,omitempty"`
	LoginAt            *string  `json:"login_at"`
	LastestLoginDevice *string  `json:"lastest_login_device,omitempty"`
	RefreshToken       *string  `json:"refresh_token,omitempty"`
	CreatedAt          *string  `json:"created_at"`
	UpdatedAt          *string  `json:"updated_at"`
}

type CreateLoginDto struct {
	UserID         string   `json:"user_id"`
	Username       string   `json:"username"`
	HashedPassword string   `json:"hashed_password"`
	Salt           string   `json:"salt,omitempty"`
	UserType       UserType `json:"user_type,omitempty"`
	DeviceID       *string  `json:"device_id"`
	IPAddress      *string  `json:"ip_address"`
	UserAgent      *string  `json:"user_agent,omitempty"`
}

type UpdateLoginDto struct {
	DeviceID       *string `json:"device_id,omitempty"`
	IPAddress      *string `json:"ip_address,omitempty"`
	UserAgent      *string `json:"user_agent,omitempty"`
	HashedPassword *string `json:"hashed_password,omitempty"`
	RefreshToken   *string `json:"refresh_token,omitempty"`
}
