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
}

func TestVectorParseEmpty(t *testing.T) {
	var vec pgvector.Vector
	err := vec.Parse("[]")
	if err != nil {
		t.Fatal(err)
	}
	if len(vec.Slice()) != 0 {
		t.Errorf("got %v, want empty", vec.Slice())
	}

	// Round-trip empty vector produced by String()
	empty := pgvector.NewVector([]float32{})
	err = vec.Parse(empty.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(vec.Slice()) != 0 {
		t.Errorf("got %v, want empty", vec.Slice())
	}
}

func TestVectorParseInvalid(t *testing.T) {
	for _, s := range []string{"", "[", "]", "1,2,3", "[1,2,3"} {
		var vec pgvector.Vector
		if err := vec.Parse(s); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", s)
		}
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
