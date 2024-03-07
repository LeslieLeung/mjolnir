package gorm

import "gorm.io/gorm"

// Paginate 分页
func Paginate(page, size, maxSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		}
		if maxSize > 0 && size > maxSize {
			size = maxSize
		}
		return db.Offset((page - 1) * size).Limit(size)
	}
}
