package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         string `gorm:"primaryKey; type:varchar(36); not null" `
	Email      string `gorm:"unique"`
	Password   string `gorm:"not null"`
	Role       int    `gorm:"default:1"`
	IsVerified bool   `gorm:"default:false"`
	Detail     UserDetail
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Users []User

func (u *User) BeforeCreate(*gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

type UserDetail struct {
	UserID    string `gorm:"primaryKey; type:varchar(36)"`
	Name      string
	Phone     string
	PictureID string         `gorm:"type:varchar(36)"`
	Picture   ProfilePicture `gorm:"foreignKey:PictureID"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProfilePicture struct {
	ID        string `gorm:"primaryKey; type:varchar(36); not null" `
	Url       string
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
