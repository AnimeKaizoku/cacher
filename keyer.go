package cacher

import (
	"fmt"
	"strings"
)

const keySep = '.'

// PolyKeyer is a special type of struct which is used when
// you want to set more than 1 parameter as a key to some value.
// It creates a single key of type string from multiple values,
// and looks like: "primaryKey.extraKey1.extraKey2.extraKeyn"
// where dot (.) works as a key separator, primaryKey is the
// 1st key and set while making a new PolyKeyer, while extraKeys
// are provided while creating keys from a polykeyer.
//
// Note: You should use Duplet in case there is only one extra
// parameter.
type PolyKeyer struct {
	primaryKey string
	extraKeys  int
	sep        rune
}

// This function creates a new PolyKeyer instance with the provided
// primary key and number of extra keys.
// Eg: If we pass "chat" to primaryKey and 2 to numExtraKeys then
// this function will create a new PolyKeyer which would create
// unified key of type string with 3 keys in it.
func NewPolyKeyer(primaryKey string, numExtraKeys int) *PolyKeyer {
	return &PolyKeyer{
		primaryKey: primaryKey,
		extraKeys:  numExtraKeys,
		sep:        keySep,
	}
}

// It creates a new unity key of type string with 1st key as
// the primary key of current PolyKeyer and rest of the keys
// in the same order as they were passed as an argument.
func (k *PolyKeyer) New(extraKeys ...string) string {
	if len(extraKeys) == 0 || (len(extraKeys) != k.extraKeys) {
		panic(fmt.Sprintf("cacher.PolyKeyer.New: invalid amound of extra keys to PolyKeyer[%s]", k.primaryKey))
	}
	n := len(extraKeys) + len(k.primaryKey)
	for i := 0; i < len(extraKeys); i++ {
		n += len(fmt.Sprint(extraKeys[i]))
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(k.primaryKey)
	for _, s := range extraKeys {
		b.WriteRune(k.sep)
		b.WriteString(s)
	}
	return b.String()
}

// Dupley is a special type of PolyKeyer which is used when
// you want to set exactly 2 parameters as a key to some value.
// It creates a single key of type string from 2 values, and
// looks like: "primaryKey.secondaryKey" where dot (.) works
// as a key separator, primaryKey is the 1st key and set while
// making a new Duplet, while secondaryKey is provided while
// creating keys from a duplet.
//
// Note: You can use PolyKeyer in case you want to use more
// than one extra parameter.
type Duplet struct {
	primaryKey string
	sep        rune
}

// This function creates a new Duplet instance with the
// provided primary key.
func NewDuplet(primaryKey string) *Duplet {
	return &Duplet{
		primaryKey: primaryKey,
		sep:        keySep,
	}
}

// It creates a new unity key of type string with 1st key as
// the primary key of current Duplet and secondary key as the
// one passed which was passed as an argument to this function.
func (k *Duplet) New(key any) string {
	return fmt.Sprintf("%s%c%v", k.primaryKey, k.sep, key)
}
