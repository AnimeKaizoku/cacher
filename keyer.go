package cacher

import (
	"fmt"
	"strings"
)

const defaultKeySep = '.'

type PolyKeyer struct {
	primaryKey string
	extraKeys  int
	sep        rune
}

func NewPolyKeyer(primaryKey string, numExtraKeys int) *PolyKeyer {
	return &PolyKeyer{
		primaryKey: primaryKey,
		extraKeys:  numExtraKeys,
		sep:        defaultKeySep,
	}
}

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

type Duplet struct {
	primaryKey string
	sep        rune
}

func NewDuplet(primaryKey string) *Duplet {
	return &Duplet{
		primaryKey: primaryKey,
		sep:        defaultKeySep,
	}
}

func (k *Duplet) New(key any) string {
	return fmt.Sprintf("%s%c%v", k.primaryKey, k.sep, key)
}
