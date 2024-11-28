package configs

import "gorm.io/gorm"

// DataInstance ...
type DataInstance struct {
	DB *gorm.DB
}

// NewDataInstance ...
func NewDataInstance(DB *gorm.DB) *DataInstance {
	return &DataInstance{DB: DB}
}
