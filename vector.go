package pgvector

import (
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
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
	buf := make([]byte, 0, 2+16*len(v.vec))
	buf = append(buf, '[')

	for i := 0; i < len(v.vec); i++ {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendFloat(buf, float64(v.vec[i]), 'f', -1, 32)
	}

	buf = append(buf, ']')
	return string(buf)
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

// EncodeBinary encodes a binary representation of the vector.
func (v Vector) EncodeBinary(buf []byte) (newBuf []byte, err error) {
	dim := len(v.vec)
	buf = slices.Grow(buf, 4+4*dim)
	buf = binary.BigEndian.AppendUint16(buf, uint16(dim))
	buf = binary.BigEndian.AppendUint16(buf, 0)
	for _, v := range v.vec {
		buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(v))
	}
	return buf, nil
}

// DecodeBinary decodes a binary representation of a vector.
func (v *Vector) DecodeBinary(buf []byte) error {
	dim := int(binary.BigEndian.Uint16(buf[0:2]))
	unused := binary.BigEndian.Uint16(buf[2:4])
	if unused != 0 {
		return fmt.Errorf("expected unused to be 0")
	}

	v.vec = make([]float32, 0, dim)
	offset := 4
	for i := 0; i < dim; i++ {
		v.vec = append(v.vec, math.Float32frombits(binary.BigEndian.Uint32(buf[offset:offset+4])))
		offset += 4
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

// statically assert that Vector implements json.Marshaler.
var _ json.Marshaler = (*Vector)(nil)

// MarshalJSON implements the json.Marshaler interface.
func (v Vector) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.vec)
}

// statically assert that Vector implements json.Unmarshaler.
var _ json.Unmarshaler = (*Vector)(nil)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vector) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.vec)
}
