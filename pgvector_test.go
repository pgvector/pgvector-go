package pgvector

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestNewVector(t *testing.T) {
	v := NewVector([]float32{1, 2, 3})
	v.Slice()
}

func TestVector_Slice(t *testing.T) {
	type fields struct {
		vec []float32
	}
	tests := []struct {
		name   string
		fields fields
		want   []float32
	}{
		{"test", fields{[]float32{1, 2, 3}}, []float32{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				vec: tt.fields.vec,
			}
			if got := v.Slice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector.Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_String(t *testing.T) {
	type fields struct {
		vec []float32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test", fields{[]float32{1, 2, 3}}, "[1,2,3]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				vec: tt.fields.vec,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Vector.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Parse(t *testing.T) {
	type fields struct {
		vec []float32
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"ok", fields{[]float32{1, 2, 3}}, args{"[1,2,3]"}, false},
		{"err", fields{[]float32{1, 2, 3}}, args{"[1,2,3"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{
				vec: tt.fields.vec,
			}
			if err := v.Parse(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vector.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVector_Scan(t *testing.T) {
	type fields struct {
		vec []float32
	}
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"string", fields{[]float32{1, 2, 3}}, args{"[1,2,3]"}, false},
		{"[]byte", fields{[]float32{1, 2, 3}}, args{[]byte("[1,2,3]")}, false},
		{"unsupported", fields{[]float32{1, 2, 3}}, args{1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{
				vec: tt.fields.vec,
			}
			if err := v.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Vector.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVector_Value(t *testing.T) {
	type fields struct {
		vec []float32
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		{"ok", fields{[]float32{1, 2, 3}}, "[1,2,3]", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				vec: tt.fields.vec,
			}
			got, err := v.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
