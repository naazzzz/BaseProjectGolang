package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required,max=191"`
	Password string `json:"password" validate:"required,min=3"`
}

type DataLoginResponse struct {
	Token     string `json:"tokendmn"`
	TokenType string `json:"token_type"`
	ExpiresAt string `json:"expires_at"`
}

type TokenInfo struct {
	Token     string `json:"tokendmn"`
	TokenType string `json:"token_type"`
	ExpiresAt string `json:"expires_at"`
}
