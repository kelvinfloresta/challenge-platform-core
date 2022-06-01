package database

import (
	"math"

	"gorm.io/gorm"
)

type PaginateInput struct {
	ActualPage int
	PageSize   int
}

func Paginate(db *gorm.DB, paginate PaginateInput) {
	offset := paginate.ActualPage * paginate.PageSize
	db.Limit(paginate.PageSize).Offset(offset)
}

func CalcMaxPages(count int64, pageSize int) int {
	total := float64(count) / float64(pageSize)
	return int(math.Ceil(total))
}
