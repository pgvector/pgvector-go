package pgvector_test

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	err = vec.Parse("[1, 2, 3]")
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
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
