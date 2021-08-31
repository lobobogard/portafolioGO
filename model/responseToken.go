package model

type ResponseToken struct {
	Token        string `json:"token"`
	TokenRefresh string `json:"tokenRefresh"`
}
