package models

import "gorm.io/gorm"

type RoomUsers struct {
	gorm.Model
	RoomId uint `gorm:"column:room_id;type:int(16);not null;" json:"room_id"`
	UserId uint `gorm:"column:user_id;type:int(16);not null;" json:"user_id"`
}

func (table *RoomUsers) TableName() string {
	return "room_users"
}
