package models

import (
	"time"

	"blog/database"
)

const (
	UserRoleNormal   = iota //0
	UserRoleUploader = iota //1
	UserRoleAdmin    = iota //2
)

type User struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	Name      *string
	Email     string `gorm:"type:text"`
	Sex       *string
	Region    *string
	CreatedAt time.Time  //帳號創建時間
	LastLogin *time.Time //最後登入時間
	Role      int        `gorm:"default:1"` //0:一般使用者 1:可上傳 2:管理員
}

type Article struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	Topic      int
	Title      string
	Author     string
	AuthorInfo string
	Time       time.Time
	Content    string `gorm:"type:text"`
	Tags       string `gorm:"type:text"`
}

type Comments struct {
	Id        uint
	ArticleId uint    //外鍵
	Article   Article `gorm:"OnDelete:SET NULL;"`
	Type      int
}

func init() {
	database.ORMModels = append(database.ORMModels, &User{}, &Article{}, &Comments{})
}
