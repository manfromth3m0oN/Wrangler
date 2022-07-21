package wrangler

import (
	"errors"
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
)

/// The incoming struct should satisfy the `Sourceable` interface which defines one method, Source
/// Source should return an map of IR structs that define the validation and data that should be validated.
/// The map key dictates which struct member the IR should map to.
/// The out struct should satisfy teh Populateable interface.
/// First the Default() method is called, which inserts any static data, then the populate function is called which
/// will place the IR into the right struct field.

var (
	ErrBadValidation = errors.New("bad validation")
	ErrBadData       = errors.New("bad data in ir")
	ErrNotPointer    = errors.New("destination is not a pointer")
)

type IRMap map[string]map[string]IR

// IR intermeditary representation of the data flowing between one struct format and another
type IR struct {
	Data       interface{}
	Validation string
}

// Validate validates the data within IR
func (ir IR) Validate() bool {
	validator, ok := govalidator.TagMap[ir.Validation]
	if !ok {
		return false
	}

	return validator(fmt.Sprintf("%v", ir.Data))
}

// Sourceable defines the methods that should be avaliable on incoming datastructures
type Sourceable interface {
	Source() IRMap
}

// Populateable defines the methods that should be avaliable on outbound datastructures
type Populateable interface {
	Default()
	Populate(IRMap) error
}

// Transform transform one source into many destinations by calling the defined interface methods
func Transform(source Sourceable, dests ...Populateable) error {
	ir := source.Source()

	for _, rm := range ir {
		for _, r := range rm {
			success := r.Validate()
			if !success {
				log.Printf("Failed to validate %v", r)
				return ErrBadValidation
			}
		}
	}

	for _, dest := range dests {
		dest.Default()

		err := dest.Populate(ir)
		if err != nil {
			return err
		}
	}

	return nil
}
