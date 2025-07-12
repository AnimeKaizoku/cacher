package cacher

import "time"

type EvictionPolicy interface {
	getEvictableValue() evictibleValue
}

type _DefaultEviction struct {
	revaluate bool
	ttl       int64
}

func (d *_DefaultEviction) getEvictableValue() evictibleValue {
	return &defaultEviction{
		revaluate: d.revaluate,
		ttl:       d.ttl,
	}
}

func DefaultEvictionPolicy(revaluate bool, ttl int64) EvictionPolicy {
	return &_DefaultEviction{
		revaluate: revaluate,
		ttl:       ttl,
	}
}

type evictibleValue interface {
	isExpired(dry bool) bool
}

type defaultEviction struct {
	expiry    int64
	revaluate bool
	ttl       int64
}

func (d *defaultEviction) isExpired(dry bool) bool {
	if d.expiry == 0 {
		return false
	}
	currTime := time.Now().Unix()
	if d.expiry <= currTime {
		return true
	}
	if dry {
		return false
	}
	if !d.revaluate {
		return false
	}
	(*d).expiry = currTime + d.ttl
	return false
}
