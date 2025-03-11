package pgvector

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// HalfVector is a wrapper for []float32 to implement sql.Scanner and driver.Valuer.
type HalfVector struct {
	vec []float32
}

// NewHalfVector creates a new HalfVector from a slice of float32.
func NewHalfVector(vec []float32) HalfVector {
	return HalfVector{vec: vec}
}

// Slice returns the underlying slice of float32.
func (v HalfVector) Slice() []float32 {
	return v.vec
}

// String returns a string representation of the half vector.
func (v HalfVector) String() string {
	// should never throw an error
	// but returning an empty string is fine if it does
	buf, _ := v.EncodeText(nil)
	return string(buf)
}

// Parse parses a string representation of a half vector.
func (v *HalfVector) Parse(s string) error {
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

// EncodeText encodes a text representation of the half vector.
func (v HalfVector) EncodeText(buf []byte) (newBuf []byte, err error) {
	buf = slices.Grow(buf, 2+16*len(v.vec))
	buf = append(buf, '[')
	for i := 0; i < len(v.vec); i++ {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendFloat(buf, float64(v.vec[i]), 'f', -1, 32)
	}
	buf = append(buf, ']')
	return buf, nil
}

// SetSlice sets the underlying slice of float32.
func (v *HalfVector) SetSlice(vec []float32) {
	v.vec = vec
}

// statically assert that HalfVector implements sql.Scanner.
var _ sql.Scanner = (*HalfVector)(nil)

// Scan implements the sql.Scanner interface.
func (v *HalfVector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// statically assert that HalfVector implements driver.Valuer.
var _ driver.Valuer = (*HalfVector)(nil)

// Value implements the driver.Valuer interface.
func (v HalfVector) Value() (driver.Value, error) {
	return v.String(), nil
}

// statically assert that HalfVector implements json.Marshaler.
var _ json.Marshaler = (*HalfVector)(nil)

// MarshalJSON implements the json.Marshaler interface.
func (v HalfVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.vec)
}

// statically assert that HalfVector implements json.Unmarshaler.
var _ json.Unmarshaler = (*HalfVector)(nil)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *HalfVector) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.vec)
}
