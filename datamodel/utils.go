package datamodel

import "gorm.io/gorm"

type IDModelIntrinsic struct {
	ID uint64 `gorm:"primary_key" json:"id"`
}

func InitDataModel(dbConn *gorm.DB) {
	_ = dbConn.AutoMigrate(&UUID{})

	_ = dbConn.AutoMigrate(&User{})
	_ = dbConn.AutoMigrate(&Article{})
}
