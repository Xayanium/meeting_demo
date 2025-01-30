package models

import (
	"gorm.io/gorm"
	"time"
)

type RoomInfo struct {
	gorm.Model
	Identity string    `gorm:"column:identity;type:varchar(64);uniqueIndex;not null;" json:"identity"` // 会议唯一标识uuid
	Name     string    `gorm:"column:name;type:varchar(64);not null;" json:"name"`
	BeginAt  time.Time `gorm:"column:begin_at;type:varchar(64);not null;" json:"begin_at"`
	EndAt    time.Time `gorm:"column:end_at;type:varchar(64);not null;" json:"end_at"`
	UserId   uint      `gorm:"column:user_id;type:int(16);not null;" json:"user_id"` // 创建人的id
}

func (table *RoomInfo) TableName() string {
	return "room_info"
}
