package domain

import "time"

type AuditLog struct {
	ID         string
	EntityID   string
	EntityType string
	Operation  string
	UserID     string
	Changes    []byte // JSON verisi
	CreatedAt  time.Time
}
