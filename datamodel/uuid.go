package datamodel

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type UUID struct {
	IDModelIntrinsic

	UUIDModel
	TimestampModel
}

func (m *UUID) Delete(dbConn *gorm.DB) error {
	return dbConn.Model(m).Where("uuid=?", m.UUID).Delete(m).Error
}

// TableName implements schema.Tabler interface
func (m *UUID) TableName() string {
	return "uuid_tab"
}

type UUIDModel struct {
	UUID string `gorm:"column:uuid;unique" json:"uuid" msgpack:"uuid"`
}

func (m *UUIDModel) GetUUID(dbConn *gorm.DB) (uuid *UUID, err error) {
	retryLimit := 5

	uuid = &UUID{}
	for count := 0; count < retryLimit; count++ {
		uuidStr := randUUIDByTimestamp()
		uuid.UUID = uuidStr
		if dbErr := dbConn.Model(uuid).Save(uuid).Error; dbErr == nil {
			return uuid, nil
		}
	}

	return nil, fmt.Errorf("could not generate valid UUID in %d retries", retryLimit)
}

func randUUIDByTimestamp() string {

	// generate random number by timestamp
	// [0 ... 27](28 bit) current timestamp, [28...63](36 bits) random number
	// (a full timestamp takes 32 bits, we emit the first 4 bits because they hardly change)
	uuidUint64 := ((uint64(time.Now().Unix()) & 0xfffffff) << 4) | (rand.Uint64() >> 28) //nolint:gosec, gomnd

	// UUID in HEX form will be a 16-char string
	return fmt.Sprintf("%016x", uuidUint64)
}
