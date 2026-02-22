package response

import "time"

type AuthLoginResponse struct {
	AccessToken  TokenResponse `json:"access_token"`
	RefreshToken TokenResponse `json:"refresh_token"`
}

type TokenResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}
