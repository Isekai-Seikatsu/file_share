package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identity struct {
	UUID       uuid.UUID `gorm:"primarykey"`
	UserAgent  string
	IP         string
	LastOnline time.Time
	CreatedAt  time.Time
}

func (i *Identity) BeforeCreate(tx *gorm.DB) error {
	i.UUID = uuid.New()
	return nil
}

func DefaultUUID() uuid.UUID {
	return uuid.New()
}

func (i *Identity) AfterFind(tx *gorm.DB) error {
	err := tx.Model(i).Update("last_online", time.Now()).Error
	return err
}