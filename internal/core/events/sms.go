package events

type SendOTPEvent struct {
	ID         string            `json:"id"`
	Recipient  string            `json:"recipient"`
	UserID     string            `json:"user_id,omitempty"`
	TemplateID string            `json:"template_id,omitempty"`
	Params     map[string]string `json:"params,omitempty"`
}
