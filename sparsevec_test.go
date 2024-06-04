package pgvector_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestSparseVectorSlice(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 0, 2, 0, 3, 0}) {
		t.Errorf("Bad slice")
	}
}

func TestSparseVectorString(t *testing.T) {
	vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
	if fmt.Sprint(vec) != "{1:1,3:2,5:3}/6" {
		t.Errorf("Bad string")
	}
}
