package cacher

import (
	"time"
)

type value[T any] struct {
	expiry int64
	val    T
}

func (v *value[T]) get(revaluate bool, ttl int64) (value T, expired bool) {
	if v.expiry == 0 {
		return v.val, false
	}
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
