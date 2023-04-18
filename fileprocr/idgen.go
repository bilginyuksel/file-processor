package fileprocr

import "github.com/google/uuid"

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (i *IDGenerator) Generate() string {
	return uuid.NewString()
}
