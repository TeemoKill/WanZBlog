package datamodel

import (
	"gorm.io/gorm"
)

type Article struct {
	IDModelIntrinsic
	Title       string `gorm:"column:title" json:"title"`
	CreatorUUID string `gorm:"column:creator_uuid" json:"creator_uuid"`

	Content string `gorm:"column:content;type:MEDIUMTEXT" json:"content"`

	UUIDModel
}

func (m *Article) TableName() string {
	return "article_tab"
}

func (m *Article) BeforeCreate(dbConn *gorm.DB) error {
	newUUID, err := m.GetUUID(dbConn)
	if err != nil {
		return err
	}
	m.UUID = newUUID.UUID

	return nil
}

func (m *Article) Create(dbConn *gorm.DB) error {
	return dbConn.Create(m).Error
}

func (m *Article) Save(dbConn *gorm.DB) error {
	return dbConn.Save(m).Error
}

func (m *Article) Delete(dbConn *gorm.DB) error {
	return dbConn.Delete(m).Error
}
