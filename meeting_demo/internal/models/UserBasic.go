package models

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(64);not null" json:"password"`
	Sdp      string `gorm:"column:sdp;type:text;not null" json:"sdp"`
}

func (table *UserInfo) TableName() string {
	return "user_info"
}
