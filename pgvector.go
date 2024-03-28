package pgvector

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Vector is a wrapper for []float32 to implement sql.Scanner and driver.Valuer.
type Vector struct {
	vec []float32
}

// NewVector creates a new Vector from a slice of float32.
func NewVector(vec []float32) Vector {
	return Vector{vec: vec}
}

// Slice returns the underlying slice of float32.
func (v Vector) Slice() []float32 {
	return v.vec
}

// String returns a string representation of the vector.
func (v Vector) String() string {
	if len(v.vec) == 0 {
		return "[]"
	}
	// brackets (2) + commas (len(v.vec)-1) + floats (len(v.vec))
	buf := make([]byte, 0, 2+2*len(v.vec)-1)
	buf = append(buf, '[')

	for i := 0; i < len(v.vec); i++ {
		buf = strconv.AppendFloat(buf, float64(v.vec[i]), 'f', -1, 32)
		buf = append(buf, ',')
	}

	buf[len(buf)-1] = ']'
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// Parse parses a string representation of a vector.
func (v *Vector) Parse(s string) error {
	sp := strings.Split(s[1:len(s)-1], ",")
	v.vec = make([]float32, 0, len(sp))
	for i := 0; i < len(sp); i++ {
		n, err := strconv.ParseFloat(sp[i], 32)
		if err != nil {
			return err
		}
		v.vec = append(v.vec, float32(n))
	}
	return nil
}

// statically assert that Vector implements sql.Scanner.
var _ sql.Scanner = (*Vector)(nil)

// Scan implements the sql.Scanner interface.
func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// statically assert that Vector implements driver.Valuer.
var _ driver.Valuer = (*Vector)(nil)

// Value implements the driver.Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	return v.String(), nil
}
