package pgvector_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestHalfVectorMarshal(t *testing.T) {
	vec := pgvector.NewHalfVector([]float32{1, 2, 3})
	data, err := json.Marshal(vec)
	if err != nil {
		panic(err)
	}
	if string(data) != `[1,2,3]` {
		t.Errorf("Bad marshal")
	}
}

func TestHalfVectorUnmarshal(t *testing.T) {
	var vec pgvector.HalfVector
	err := json.Unmarshal([]byte(`[1,2,3]`), &vec)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(vec.Slice(), []float32{1, 2, 3}) {
		t.Errorf("Bad unmarshal")
	}
}
