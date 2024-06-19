package generator

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/wisle25/be-template/applications/generator"
)

// UUIDGenerator implements IdGenerator using UUID
type UUIDGenerator struct /* implements IdGenerator */ {

}

func NewUUIDGenerator() generator.IdGenerator {
	return &UUIDGenerator{}
}

func (generator *UUIDGenerator) Generate() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("id_generator_err: generate: %v", err))
	}

	return id.String()
}
