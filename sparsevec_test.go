package pgvector_test

import (
	"fmt"
	"reflect"
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
}
