package pgvector_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestHalfVectorSlice(t *testing.T) {
	vec := pgvector.NewHalfVector([]float32{1, 2, 3})
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestHalfVectorString(t *testing.T) {
	vec := pgvector.NewHalfVector([]float32{1, 2, 3})
	if fmt.Sprint(vec) != "[1,2,3]" {
		t.Error()
	}
}

func TestHalfVectorParse(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("[1,2,3]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}

func TestHalfVectorParseZeroDimensions(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("[]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{}) {
		t.Error()
	}
}

func TestHalfVectorParseZeroLength(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("")
	if err == nil || err.Error() != "malformed halfvec literal" {
		t.Error()
	}
}

// TODO error in 0.5.0
func TestHalfVectorParseOpenBracket(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("{1]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1}) {
		t.Error()
	}
}

// TODO error in 0.5.0
func TestHalfVectorParseCloseBracket(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("[1}")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1}) {
		t.Error()
	}
}

func TestHalfVectorParseNonNumber(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("[a]")
	if err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Error()
	}
}

func TestHalfVectorParseOutOfRange(t *testing.T) {
	var vec pgvector.HalfVector
	err := vec.Parse("[4e38]")
	if err == nil || !errors.Is(err, strconv.ErrRange) {
		t.Error()
	}
}

func TestHalfVectorEncodeText(t *testing.T) {
	vec := pgvector.NewHalfVector([]float32{1, 2, 3})
	buf := make([]byte, 0)
	newBuf, err := vec.EncodeText(buf)
	if err != nil {
		panic(err)
	}
	if string(newBuf) != "[1,2,3]" {
		t.Error()
	}
}

func TestHalfVectorMarshal(t *testing.T) {
	vec := pgvector.NewHalfVector([]float32{1, 2, 3})
	data, err := json.Marshal(vec)
	if err != nil {
		panic(err)
	}
	if string(data) != "[1,2,3]" {
		t.Error()
	}
}

func TestHalfVectorUnmarshal(t *testing.T) {
	var vec pgvector.HalfVector
	err := json.Unmarshal([]byte("[1,2,3]"), &vec)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Error()
	}
}
