package wrangler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSource struct {
	X string
	Y int
	Z string
}

func (s TestSource) Source() IRMap {
        return map[string]map[string]IR {
          "X": {"val": {Data: s.X, Validation: "ascii"}},
          "Y": {"val": {Data: s.Y, Validation: "numeric"}},
          "Z": {"val": {Data: s.Z, Validation: "ascii"}},
        }
}

type TestDestX struct {
	ID   string
	Type string
	Data string
}

func (d *TestDestX) Default() {
	d.ID = uuid.New().String()
	d.Type = "destinationX"
}

func (d *TestDestX) Populate(ir IRMap) error {
	data, ok := ir["X"]["val"].Data.(string)
	if !ok {
		return ErrBadData
	}

	d.Data = data

	return nil
}

type TestDestY struct {
	ID   string
	Type string
	Data int
}

func (d *TestDestY) Default() {
	d.ID = uuid.New().String()
	d.Type = "destinationY"
}

func (d *TestDestY) Populate(ir IRMap) error {
	data, ok := ir["Y"]["val"].Data.(int)
	if !ok {
		return ErrBadData
	}

	d.Data = data

	return nil
}

type TestDestZ struct {
	ID   string
	Type string
	Data string
}

func (d *TestDestZ) Default() {
	d.ID = uuid.New().String()
	d.Type = "destinationY"
}

func (d *TestDestZ) Populate(ir IRMap) error {
	data, ok := ir["Z"]["val"].Data.(string)
	if !ok {
		return ErrBadData
	}

	d.Data = data

	return nil
}

func TestTransform(t *testing.T) {
	source := TestSource{
		X: "Test data x",
		Y: 42,
		Z: "Test data y",
	}

	var destX TestDestX
	var destY TestDestY
	var destZ TestDestZ
	err := Transform(source, &destX, &destY, &destZ)
        t.Log(destX)
        t.Log(destY)
        t.Log(destZ)
	assert.NoError(t, err)
	assert.Equal(t, source.X, destX.Data)
	assert.Equal(t, source.Y, destY.Data)
	assert.Equal(t, source.Z, destZ.Data)
}
