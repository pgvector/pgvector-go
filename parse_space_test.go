package pgvector

import (
	"reflect"
	"testing"
)

func TestParseTrimSpace(t *testing.T) {
	var v Vector
	if err := v.Parse(" [1, 2, 3] "); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v.Slice(), []float32{1, 2, 3}) {
		t.Fatalf("%v", v.Slice())
	}
	var h HalfVector
	if err := h.Parse(" [4, 5] "); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(h.Slice(), []float32{4, 5}) {
		t.Fatalf("%v", h.Slice())
	}
}
