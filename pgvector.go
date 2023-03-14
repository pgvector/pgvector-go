package pgvector

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Vector struct {
	vec []float32
}

func NewVector(vec []float32) Vector {
	return Vector{vec: vec}
}

func (v Vector) Slice() []float32 {
	return v.vec
}

func (v Vector) String() string {
	var buf strings.Builder
	buf.WriteString("[")

	for i := 0; i < len(v.vec); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatFloat(float64(v.vec[i]), 'f', -1, 32))
	}

	buf.WriteString("]")
	return buf.String()
}

func (v *Vector) Parse(s string) error {
	v.vec = make([]float32, 0)
	sp := strings.Split(s[1:len(s)-1], ",")
	for i := 0; i < len(sp); i++ {
		n, err := strconv.ParseFloat(sp[i], 32)
		if err != nil {
			return err
		}
		v.vec = append(v.vec, float32(n))
	}
	return nil
}

var _ sql.Scanner = (*Vector)(nil)

func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

var _ driver.Valuer = (*Vector)(nil)

func (v Vector) Value() (driver.Value, error) {
	return v.String(), nil
}
