package helpers

import (
	"sudoku-daily-api/src/domain/helpers"

	"github.com/google/uuid"
)

type (
	uuidHelper struct{}
)

func NewUUIDHelper() helpers.UUIDHelper {
	return &uuidHelper{}
}

func (u *uuidHelper) NewUUID() string {
	return uuid.NewString()
}
