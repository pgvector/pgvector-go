package pgvector

import (
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// SparseVector is a wrapper to implement sql.Scanner and driver.Valuer.
type SparseVector struct {
	dim     int32
	indices []int32
	values  []float32
}

// NewSparseVector creates a new SparseVector from a slice of float32.
func NewSparseVector(vec []float32) SparseVector {
	dim := int32(len(vec))
	indices := make([]int32, 0)
	values := make([]float32, 0)
	for i := 0; i < len(vec); i++ {
		if vec[i] != 0 {
			indices = append(indices, int32(i))
			values = append(values, vec[i])
		}
	}
	return SparseVector{dim: dim, indices: indices, values: values}
}

// NewSparseVectorFromMap creates a new SparseVector from a map of non-zero elements.
func NewSparseVectorFromMap(elements map[int32]float32, dim int32) SparseVector {
	indices := make([]int32, 0, len(elements))
	values := make([]float32, 0, len(elements))
	for k, v := range elements {
		if v != 0 {
			indices = append(indices, k)
		}
	}
	slices.Sort(indices)
	for _, k := range indices {
		values = append(values, elements[k])
	}
	v := SparseVector{dim: dim, indices: indices, values: values}
	err := v.validate()
	if err != nil {
		// TODO return error in 0.5.0
		panic(err)
	}
	return v
}

// Dimensions returns the number of dimensions.
func (v SparseVector) Dimensions() int32 {
	return v.dim
}

// Indices returns the non-zero indices.
func (v SparseVector) Indices() []int32 {
	return v.indices
}

// Values returns the non-zero values.
func (v SparseVector) Values() []float32 {
	return v.values
}

// Slice returns a slice of float32.
func (v SparseVector) Slice() []float32 {
	vec := make([]float32, v.dim)
	for i := 0; i < len(v.indices); i++ {
		vec[v.indices[i]] = v.values[i]
	}
	return vec
}

// String returns a string representation of the sparse vector.
func (v SparseVector) String() string {
	buf := make([]byte, 0, 13+27*len(v.indices))
	buf = append(buf, '{')

	for i := 0; i < len(v.indices); i++ {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendInt(buf, int64(v.indices[i])+1, 10)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, float64(v.values[i]), 'f', -1, 32)
	}

	buf = append(buf, '}')
	buf = append(buf, '/')
	buf = strconv.AppendInt(buf, int64(v.dim), 10)
	return string(buf)
}

// Parse parses a string representation of a sparse vector.
func (v *SparseVector) Parse(s string) error {
	sp := strings.SplitN(s, "/", 2)
	if len(sp) != 2 {
		return fmt.Errorf("malformed sparsevec literal")
	}

	dim, err := strconv.ParseInt(sp[1], 10, 32)
	if err != nil {
		return err
	}

	// TODO check brackets in 0.5.0
	if len(sp[0]) < 2 {
		return fmt.Errorf("malformed sparsevec literal")
	}

	elements := []string{}
	if len(sp[0]) > 2 {
		elements = strings.Split(sp[0][1:len(sp[0])-1], ",")
	}

	v.dim = int32(dim)
	v.indices = make([]int32, 0, len(elements))
	v.values = make([]float32, 0, len(elements))

	for i := 0; i < len(elements); i++ {
		ep := strings.SplitN(elements[i], ":", 2)
		if len(ep) != 2 {
			return fmt.Errorf("malformed sparsevec literal")
		}

		n, err := strconv.ParseInt(ep[0], 10, 32)
		if err != nil {
			return err
		}
		v.indices = append(v.indices, int32(n-1))

		n2, err := strconv.ParseFloat(ep[1], 32)
		if err != nil {
			return err
		}
		v.values = append(v.values, float32(n2))
	}

	return v.validate()
}

// EncodeBinary encodes a binary representation of the sparse vector.
func (v SparseVector) EncodeBinary(buf []byte) (newBuf []byte, err error) {
	nnz := len(v.indices)
	buf = slices.Grow(buf, 12+8*nnz)
	buf = binary.BigEndian.AppendUint32(buf, uint32(v.dim))
	buf = binary.BigEndian.AppendUint32(buf, uint32(nnz))
	buf = binary.BigEndian.AppendUint32(buf, 0)
	for _, v := range v.indices {
		buf = binary.BigEndian.AppendUint32(buf, uint32(v))
	}
	for _, v := range v.values {
		buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(v))
	}
	return buf, nil
}

// DecodeBinary decodes a binary representation of a sparse vector.
func (v *SparseVector) DecodeBinary(buf []byte) error {
	if len(buf) < 12 {
		return fmt.Errorf("invalid length")
	}

	dim := int32(binary.BigEndian.Uint32(buf[0:4]))
	nnz := int(binary.BigEndian.Uint32(buf[4:8]))
	if nnz < 0 {
		return fmt.Errorf("sparsevec cannot have negative number of elements")
	}

	unused := binary.BigEndian.Uint32(buf[8:12])
	if unused != 0 {
		return fmt.Errorf("expected unused to be 0")
	}

	if (len(buf)-12)/8 != nnz || (len(buf)-12)%8 != 0 {
		return fmt.Errorf("invalid length")
	}

	v.dim = dim
	v.indices = make([]int32, 0, nnz)
	v.values = make([]float32, 0, nnz)
	offset := 12

	for i := 0; i < nnz; i++ {
		v.indices = append(v.indices, int32(binary.BigEndian.Uint32(buf[offset:offset+4])))
		offset += 4
	}

	for i := 0; i < nnz; i++ {
		v.values = append(v.values, math.Float32frombits(binary.BigEndian.Uint32(buf[offset:offset+4])))
		offset += 4
	}

	return v.validate()
}

func (v *SparseVector) validate() error {
	if v.dim < 0 {
		return fmt.Errorf("sparsevec cannot have negative dimensions")
	}

	for _, index := range v.indices {
		if index < 0 || index >= v.dim {
			return fmt.Errorf("sparsevec index out of bounds")
		}
	}

	return nil
}

// statically assert that SparseVector implements sql.Scanner.
var _ sql.Scanner = (*SparseVector)(nil)

// Scan implements the sql.Scanner interface.
func (v *SparseVector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// statically assert that SparseVector implements driver.Valuer.
var _ driver.Valuer = (*SparseVector)(nil)

// Value implements the driver.Valuer interface.
func (v SparseVector) Value() (driver.Value, error) {
	return v.String(), nil
}
