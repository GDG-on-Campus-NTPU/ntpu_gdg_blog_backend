package models

import (
	"time"

	"blog/database"

	"gorm.io/datatypes"
)

const (
	UserRoleNormal   = iota //0
	UserRoleUploader = iota //1
	UserRoleAdmin    = iota //2
)

type User struct {
	Id           uint `gorm:"primaryKey;autoIncrement"`
	Name         *string
	Email        string `gorm:"type:text"`
	Sex          *string
	Region       *string
	CreatedAt    time.Time  //帳號創建時間
	LastLogin    *time.Time //最後登入時間
	Role         int        `gorm:"default:0"` //0:一般使用者 1:可上傳 2:管理員
	Article      []Article  `gorm:"foreignKey:UserId;constraint:OnDelete:SET NULL;"`
	ProfilePhoto *string    `gorm:"type:text"` //照片
	Description  *string    `gorm:"type:text"` //個人簡介
	Avatar       *string    `gorm:"type:text"` //頭像
	Major        *string    `gorm:"type:text"` //科系
}

type Article struct {
	Id          uint `gorm:"primaryKey;autoIncrement"`
	Topic       int
	Title       string
	Author      string
	AuthorImage string
	AuthorInfo  string
	Time        time.Time
	Content     string         `gorm:"type:text"`
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	UserId      uint
	Type        int // 文章類型 {1: 技術 2:回顧}
	Description string
}

type Comments struct {
	Id        uint
	ArticleId uint    //外鍵
	Article   Article `gorm:"OnDelete:SET NULL;"`
	Type      int
}

// title, thumbnail, tag, description, images, start date, end date, members(profile photo, name, major, description, links)
type Project struct {
	Id          uint `gorm:"primaryKey;autoIncrement"`
	Title       string
	Thumbnail   string
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	Description string
	Images      datatypes.JSON `gorm:"type:jsonb"`
	StartDate   time.Time
	EndDate     time.Time
	Members     []User `gorm:"many2many:project_members;OnDelete:CASCADE;"`
	UserId      uint   `gorm:"many2many:project_uploader;"`
}

// thumbnail, title, date, description, tags
type Activity struct {
	Id          uint `gorm:"primaryKey;autoIncrement"`
	Thumbnail   string
	Title       string
	Date        time.Time
	Description string
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	UserId      uint
}

func init() {
	database.ORMModels = append(database.ORMModels, &User{}, &Article{}, &Comments{}, &Project{}, &Activity{})
}
