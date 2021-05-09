package models

import (
	"time"

	"github.com/google/uuid"
)

type Bucket struct {
	ID           uint      `gorm:"primaryKey"`
	IdentityUUID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Expiry   time.Time
	CreatedAt    time.Time
}

func NewBucket(identityUUID uuid.UUID, expire time.Duration) *Bucket {
	return &Bucket{IdentityUUID: identityUUID, Expiry: time.Now().Add(expire)}
}