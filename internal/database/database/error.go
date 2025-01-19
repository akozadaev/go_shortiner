package database

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNotFound    = errors.New("record not found")
	ErrKeyConflict = errors.New("key conflict")
)

// IsRecordNotFoundErr returns true if err is gorm.ErrRecordNotFound or ErrNotFound
func IsRecordNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound)
}

// IsKeyConflictErr returns true if err is ErrKeyConflict or MySQLError with 1062 code number
func IsKeyConflictErr(err error) bool {
	if errors.Is(err, ErrKeyConflict) || errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	return false
}
