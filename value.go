package cacher

import (
	"time"
)

type Value[T any] struct {
	expiry int64
	val    T
}

func (v *Value[T]) Get(revaluate bool, ttl int64) (value T, expired bool) {
	currTime := time.Now().Unix()
	if v.expiry <= currTime {
		expired = true
		return
	}
	value = v.val
	if !revaluate {
		return
	}
	(*v).expiry = currTime + ttl
	return
}
