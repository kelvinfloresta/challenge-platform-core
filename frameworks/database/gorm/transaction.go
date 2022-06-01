package database

import "gorm.io/gorm"

type Transaction interface {
	Transaction(fc func(tx *gorm.DB) error) (err error)
}
