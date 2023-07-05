package datamodel

import (
	"gorm.io/gorm"
)

type User struct {
	IDModelIntrinsic
	Email    string `gorm:"column:email;unique" json:"email"`
	Password string `gorm:"column:password" json:"password"`
	Username string `gorm:"column:username" json:"username"`
	Bio      string `gorm:"column:bio;type:VARCHAR(300)" json:"bio"`

	UUIDModel
	TimestampModel
}

func (u *User) TableName() string {
	return "user_tab"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if uuid, err := u.GetUUID(db); err != nil {
		return err
	} else {
		u.UUID = uuid.UUID
		return nil
	}
}

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) LoadByUUID(db *gorm.DB) error {
	return db.Model(u).Where("uuid=?", u.UUID).First(u).Error
}

func (u *User) LoadByEmail(db *gorm.DB) error {
	return db.Model(u).Where("email=?", u.Email).First(u).Error
}
