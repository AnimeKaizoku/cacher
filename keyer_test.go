package cacher

import (
	"testing"
)

func TestPolyKeyer_New(t *testing.T) {
	type fields struct {
		primaryKey string
		extraKeys  int
	}
	tests := []struct {
		name      string
		fields    fields
		extraKeys []string
		want      string
	}{
		{
			name:      "dupley key system",
			fields:    fields{"pkey", 1},
			extraKeys: []string{"100"},
			want:      "pkey.100",
		},
		{
			name:      "huge key system",
			fields:    fields{"pkey", 6},
			extraKeys: []string{"100", "dde", "nice", "wow", "kek", "brah"},
			want:      "pkey.100.dde.nice.wow.kek.brah",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &PolyKeyer{
				primaryKey: tt.fields.primaryKey,
				extraKeys:  tt.fields.extraKeys,
				sep:        '.',
			}
			if got := k.New(tt.extraKeys...); got != tt.want {
				t.Errorf("PolyKeyer.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuplet_New(t *testing.T) {
	tests := []struct {
		name       string
		primaryKey string
		key        any
		want       string
	}{
		{"string.int", "pkey", 200, "pkey.200"},
		{"string.string", "pkey", "world", "pkey.world"},
		{"string.struct", "pkey", struct {
			a int
			b string
		}{10, "nice"}, "pkey.{10 nice}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Duplet{
				primaryKey: tt.primaryKey,
				sep:        '.',
			}
			if got := k.New(tt.key); got != tt.want {
				t.Errorf("Duplet.New() = %v, want %v", got, tt.want)
			}
		})
	}
}
