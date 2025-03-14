package pgvector_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/cookieai-jar/pgvector-go"
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
