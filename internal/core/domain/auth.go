package domain

type Auth struct {
	User
	AccessToken string `json:"access_token"`
}

type AuthDTO struct {
	AccessToken string `json:"access_token"`
}

type AuthRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=admin enduser drone"`
}

func (a *Auth) ToDTO() *AuthDTO {
	return &AuthDTO{
		AccessToken: a.AccessToken,
	}
}
