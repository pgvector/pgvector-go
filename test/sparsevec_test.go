package pgvector_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestNewSparseVector(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Error()
	}
}

func TestNewSparseVectorFromMap(t *testing.T) {
	vec := pgvector.NewSparseVectorFromMap(map[int32]float32{2: 2, 4: 3, 0: 1, 3: 0}, 6)
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Error()
	}
}

func TestSparseVectorDimensions(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if vec.Dimensions() != 6 {
		t.Error()
	}
}

func TestSparseVectorIndices(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if !reflect.DeepEqual(vec.Indices(), []int32{0, 2, 4}) {
		t.Error()
	}
}

func TestSparseVectorValues(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if !reflect.DeepEqual(vec.Values(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestSparseVectorSlice(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Error()
	}
}

func TestSparseVectorString(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if fmt.Sprint(vec) != "{1:1,3:2,5:3}/6" {
		t.Error()
	}
}

func TestSparseVectorFromMapString(t *testing.T) {
	vec := pgvector.NewSparseVectorFromMap(map[int32]float32{2: 2, 4: 3, 0: 1, 3: 0}, 6)
	if fmt.Sprint(vec) != "{1:1,3:2,5:3}/6" {
		t.Error()
	}

	// TODO test
	// vec = pgvector.NewSparseVectorFromMap(map[int32]float32{}, -1)
	// vec = pgvector.NewSparseVectorFromMap(map[int32]float32{1: 1}, 0)
}

func TestSparseVectorParse(t *testing.T) {
	var vec pgvector.SparseVector
	err := vec.Parse("{1:1,3:2,5:3}/6")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Error()
	}

	err = vec.Parse("{}/0")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{}) {
		t.Error()
	}

	err = vec.Parse("")
	if err == nil || err.Error() != "malformed sparsevec literal" {
		t.Error()
	}

	err = vec.Parse("/6")
	if err == nil || err.Error() != "malformed sparsevec literal" {
		t.Error()
	}

	err = vec.Parse("{1}/6")
	if err == nil || err.Error() != "malformed sparsevec literal" {
		t.Error()
	}

	err = vec.Parse("{}/-1")
	if err == nil || err.Error() != "sparsevec cannot have negative dimensions" {
		t.Error()
	}

	err = vec.Parse("{1:1}/0")
	if err == nil || err.Error() != "sparsevec index out of bounds" {
		t.Error()
	}

	err = vec.Parse("{}/a")
	if err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Error()
	}

	err = vec.Parse("{a:1}/6")
	if err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Error()
	}

	err = vec.Parse("{1:a}/6")
	if err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Error()
	}

	err = vec.Parse("{1:4e38}/6")
	if err == nil || !errors.Is(err, strconv.ErrRange) {
		t.Error()
	}
}

func TestSparseVectorEncodeBinary(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	buf := make([]byte, 0)
	got, err := vec.EncodeBinary(buf)
	if err != nil {
		panic(err)
	}
	want := make([]byte, 0)
	want = binary.BigEndian.AppendUint32(want, 6)
	want = binary.BigEndian.AppendUint32(want, 3)
	want = binary.BigEndian.AppendUint32(want, 0)
	want = binary.BigEndian.AppendUint32(want, 0)
	want = binary.BigEndian.AppendUint32(want, 2)
	want = binary.BigEndian.AppendUint32(want, 4)
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(1.0))
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(2.0))
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(3.0))
	if !reflect.DeepEqual(got, want) {
		t.Error()
	}
}

func TestSparseVectorDecodeBinary(t *testing.T) {
	var vec pgvector.SparseVector
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint32(buf, 6)
	buf = binary.BigEndian.AppendUint32(buf, 3)
	buf = binary.BigEndian.AppendUint32(buf, 0)
	buf = binary.BigEndian.AppendUint32(buf, 0)
	buf = binary.BigEndian.AppendUint32(buf, 2)
	buf = binary.BigEndian.AppendUint32(buf, 4)
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(1.0))
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(2.0))
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(3.0))
	err := vec.DecodeBinary(buf)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Error()
	}
}

func TestSparseVectorDecodeBinaryZeroLength(t *testing.T) {
	var vec pgvector.SparseVector
	buf := make([]byte, 0)
	err := vec.DecodeBinary(buf)
	if err == nil || err.Error() != "invalid length" {
		t.Error()
	}
}

func TestSparseVectorDecodeBinaryNonZeroUnused(t *testing.T) {
	var vec pgvector.SparseVector
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	err := vec.DecodeBinary(buf)
	if err == nil || err.Error() != "expected unused to be 0" {
		t.Error()
	}
}
