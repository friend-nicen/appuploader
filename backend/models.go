package backend

import (
	"time"
)

type ApiKey struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"type:text;not null" json:"name"`
	IssuerID   string    `gorm:"type:text;not null" json:"issuer_id"`
	KeyID      string    `gorm:"type:text;not null" json:"key_id"`
	PrivateKey string    `gorm:"type:text;not null" json:"private_key"`
	IsActive   bool      `gorm:"default:false" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
