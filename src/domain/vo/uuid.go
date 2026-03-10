package vo

import "github.com/google/uuid"

type UUID string

func NewUUID() UUID {
	return UUID(uuid.NewString())
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) IsValid() bool {
	_, err := uuid.Parse(string(u))
	return err == nil
}