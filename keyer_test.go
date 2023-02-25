package cacher

import (
	"fmt"
	"testing"
)

func ExamplePolyKeyer() {
	// 2 is number of extra keys we'd use.
	// total keys become 3 then
	keyer := NewPolyKeyer("chat", 2)

	fmt.Println(keyer.New("public", fmt.Sprint(100291)))

	fmt.Println(keyer.New("private", fmt.Sprint(100292)))
	// Output: chat.public.100291
	// chat.private.100292
}

func ExamplePolyKeyer_New() {
	// 2 is number of extra keys we'd use.
	// total keys become 3 then
	keyer := NewPolyKeyer("chat", 2)

	fmt.Println(keyer.New("public", fmt.Sprint(100291)))

	fmt.Println(keyer.New("private", fmt.Sprint(100292)))
	// Output: chat.public.100291
	// chat.private.100292
}

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

func ExampleDuplet() {
	// 2 is number of extra keys we'd use.
	// total keys become 3 then
	keyer := NewDuplet("chat")

	fmt.Println(keyer.New(fmt.Sprint(100291)))

	fmt.Println(keyer.New(fmt.Sprint(100292)))
	// Output: chat.100291
	// chat.100292
}

func ExampleDuplet_New() {
	// 2 is number of extra keys we'd use.
	// total keys become 3 then
	keyer := NewPolyKeyer("chat", 2)

	fmt.Println(keyer.New("public", fmt.Sprint(100291)))

	fmt.Println(keyer.New("private", fmt.Sprint(100292)))
	// Output: chat.public.100291
	// chat.private.100292
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
