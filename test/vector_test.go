package pgvector_test

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestVectorSlice(t *testing.T) {
	vec := pgvector.NewVector([]float32{1, 2, 3})
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestVectorString(t *testing.T) {
	vec := pgvector.NewVector([]float32{1, 2, 3})
	if fmt.Sprint(vec) != "[1,2,3]" {
		t.Error()
	}
}

func TestVectorParse(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[1,2,3]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestVectorParseZeroDimensions(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{}) {
		t.Error()
	}
}

func TestVectorParseZeroLength(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("")
	if err == nil || err.Error() != "malformed vector literal" {
		t.Error()
	}
}

// TODO error in 0.5.0
func TestVectorParseOpenBracket(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("{1]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1}) {
		t.Error()
	}
}

// TODO error in 0.5.0
func TestVectorParseCloseBracket(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[1}")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1}) {
		t.Error()
	}
}

func TestVectorParseNonNumber(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[a]")
	if err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Error()
	}
}

func TestVectorParseOutOfRange(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[4e38]")
	if err == nil || !errors.Is(err, strconv.ErrRange) {
		t.Error()
	}
}

func TestVectorEncodeBinary(t *testing.T) {
	vec := pgvector.NewVector([]float32{1, 2, 3})
	buf := make([]byte, 0)
	got, err := vec.EncodeBinary(buf)
	if err != nil {
		panic(err)
	}
	want := make([]byte, 0)
	want = binary.BigEndian.AppendUint16(want, 3)
	want = binary.BigEndian.AppendUint16(want, 0)
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(1.0))
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(2.0))
	want = binary.BigEndian.AppendUint32(want, math.Float32bits(3.0))
	if !reflect.DeepEqual(got, want) {
		t.Error()
	}
}

func TestVectorDecodeBinary(t *testing.T) {
	var vec pgvector.Vector
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint16(buf, 3)
	buf = binary.BigEndian.AppendUint16(buf, 0)
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(1.0))
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(2.0))
	buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(3.0))
	err := vec.DecodeBinary(buf)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestVectorDecodeBinaryZeroLength(t *testing.T) {
	var vec pgvector.Vector
	buf := make([]byte, 0)
	err := vec.DecodeBinary(buf)
	if err == nil || err.Error() != "invalid length" {
		t.Error()
	}
}

func TestVectorDecodeBinaryNonZeroUnused(t *testing.T) {
	var vec pgvector.Vector
	buf := []byte{0, 0, 0, 1}
	err := vec.DecodeBinary(buf)
	if err == nil || err.Error() != "expected unused to be 0" {
		t.Error()
	}
}

func TestVectorMarshal(t *testing.T) {
	vec := pgvector.NewVector([]float32{1, 2, 3})
	data, err := json.Marshal(vec)
	if err != nil {
		panic(err)
	}
	if string(data) != "[1,2,3]" {
		t.Error()
	}
}

func TestVectorUnmarshal(t *testing.T) {
	var vec pgvector.Vector
	err := json.Unmarshal([]byte("[1,2,3]"), &vec)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}
