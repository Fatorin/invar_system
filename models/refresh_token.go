package models

import (
	"time"
)

type RefreshToken struct {
	Model
	RoleType    uint      `json:"role_type"`
	RoleID      uint      `json:"role_id"`
	Token       string    `json:"token"`
	Expires     time.Time `json:"expires"`
	Revoked     time.Time `json:"revoked"`
	CreatedByIP string    `json:"created_by_ip"`
	RevokedByIP string    `json:"revoked_by_ip"`
}

func (refreshToken *RefreshToken) IsRevoked() bool {
	if refreshToken.Revoked.IsZero() {
		return false
	}
	return true
}

func (refreshToken *RefreshToken) SetToRevoked(ip string) {
	refreshToken.Revoked = time.Now()
	refreshToken.RevokedByIP = ip
}
