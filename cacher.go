package cacher

import (
	"sync"
	"time"
)

// Cacher is a fast, decentralised caching system, generic
// in nature and uses Go's in built mapping to store data.
// It has a plenty of features like TTL, Revaluation etc.
//
// Some of the main features are descried below:
//
// TTL (Time-To-Live): It allows us to expire a key after
// a specific time period.
// Eg: If TTL is set to 30 seconds, then each key of current
// Cacher will be expired after 30 seconds of their addition.
//
// Revaluation: This is another useful feature that allows us
// to keep keys cached as per their usage frequency.
// Working: Whenever the keys will be retrieved via (Cacher.Get)
// method, its expiry will be renewed and this will allow us to
// keep frequently used keys in the map without expiration.
type Cacher[C comparable, T any] struct {
	ttl           int64
	mutex         *sync.RWMutex
	status        status
	cacheMap      map[C]*value[T]
	revaluate     bool
	cleanInterval time.Duration
}

// This struct contains the optional arguments which can be filled
// while creating a new Cacher instance.
//
// Parameters:
//
// TimeToLive (type time.Duration):
// It allows us to expire a key after a specific time period.
// Eg: If it is set to 30 seconds, then each key of current
// Cacher will be expired after 30 seconds of their addition.
//
// CleanInterval (type time.Duration):
// It is the time of interval between two cleaner windows.
// A cleaner window is that time frame when all the expired
// keys will be deleted from our cache mapping.
// Note: It TTL is set to a finite value and no value is passed
// to CleanInterval, it'll use a default time interval of 1 day
// for clean window.
// Eg: If CleanInterval is set to 1 hour, then cleaner
// window will be run after every 1 hour, and the expired keys
// which are present in our cache map will be deleted.
//
// Revaluate (type bool):
// It allows us to keep keys cached as per their usage frequency.
// Working: Whenever the keys will be retrieved via (Cacher.Get)
// method, its expiry will be renewed and this will allow us to
// keep frequently used keys in the map without expiration.
type NewCacherOpts struct {
	TimeToLive    time.Duration
	CleanInterval time.Duration
	Revaluate     bool
}

// NewCacher is a generic function which creates a new Cacher instance.
//
// Generic parameters (for current Cacher instance):
//
// KeyT: It is the "static" type of keys of our cache.
// It accepts types which implement built-in comparable interface.
// Eg: If it is set to string, then keys will only be allowed
// as a string built-in data-type.
//
// ValueT: It is the type of values of our cache.
// It can be set to any type.
// Eg: If it is set to string, then value will only be allowed
// as a string built-in data-type.
//
// Input parameters:
//
// opts (type *NewCacherOpts):
// It contains optional parameters which you can use while creating
// a new Cacher instance.
//
// General Example:
// c := cacher.NewCacher[int, string](&cacher.NewCacherOpts{10*time.Minute, time.Hour, true})
// will create a new Cacher instance which will expire keys after 10
// minutes of their addition to the system, all the expired keys will
// be deleted from cache once in an hour. Keys will have their expiry
// revalueted on every c.Get call.
func NewCacher[KeyT comparable, ValueT any](opts *NewCacherOpts) *Cacher[KeyT, ValueT] {
	if opts == nil {
		opts = new(NewCacherOpts)
	}
	ttl := int64(opts.TimeToLive.Seconds())
	c := Cacher[KeyT, ValueT]{
		cacheMap:      make(map[KeyT]*value[ValueT]),
		mutex:         new(sync.RWMutex),
		ttl:           ttl,
		cleanInterval: opts.CleanInterval,
		revaluate:     opts.Revaluate,
	}
	if ttl != 0 {
		if c.cleanInterval == 0 {
			c.cleanInterval = time.Hour * 24
		}
		go c.cleaner()
	}
	return &c
}

// Set is used to set a new key-value pair to the current
// Cacher instance. It doesn't return anything.
func (c *Cacher[C, T]) Set(key C, val T) {
	c.setRawValue(key, c.packValue(val))
}

func (c *Cacher[C, T]) setRawValue(key C, val *value[T]) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap[key] = val
}

// Set is used to get value of the input key. It returns
// value of input key with true while returns empty value
// with false if key is not found or has expired already
//
// Note: It will renew the expiration time of the input
// key which is retrieved if revaluation mode is on for
// current Cacher instance.
func (c *Cacher[C, T]) Get(key C) (value T, ok bool) {
	rValue, ok := c.getRawValue(key)
	if !ok {
		return
	}
	val, expired := rValue.get(c.revaluate, c.ttl)
	if !expired {
		value = val
		return
	}
	ok = false
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cacheMap, key)
	return
}

// GetAll is used to return all the unexpired key-value
// pairs present in the current Cacher instance, returns
// a slice of values.
//
// Note: It doesn't renew expiration time of any key
// even if the revaluation mode is turned on for the
// current Cacher instance.
func (c *Cacher[C, T]) GetAll() []T {
	res := make([]T, c.NumKeys())
	var i = 0
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, rv := range c.cacheMap {
		v, exp := rv.get(false, 0)
		if exp {
			continue
		}
		res[i] = v
		i++
	}
	return res
}

// SegrigatorFunc takes the input as value of current key.
// Returned boolean is used for segrigation of keys for
// GetSome function.
type SegrigatorFunc[T any] func(value T) bool

// GetSome is used to get keys which satisfired a particular
// condition determined via SegrigatorFunc.
// It returns those values which satisfied the condition
// determined via SegrigatorFunc.
func (c *Cacher[C, T]) GetSome(cond SegrigatorFunc[T]) []T {
	if cond == nil {
		cond = func(T) bool { return true }
	}
	return c.getSome(cond)
}

// The inner function of GetSome
func (c *Cacher[C, T]) getSome(cond SegrigatorFunc[T]) []T {
	// we can't determine length yet due to segrigations by the
	// cond function.
	res := []T{}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, rv := range c.cacheMap {
		// No need to pass actual ttl since we ain't revaluating
		v, exp := rv.get(false, 0)
		if exp {
			continue
		}
		if !cond(v) {
			continue
		}
		res = append(res, v)
	}
	return res
}

// It returns the value of a key in the form of Value struct.
func (c *Cacher[C, T]) getRawValue(key C) (val *value[T], ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok = c.cacheMap[key]
	return
}

// It packs the value to a a struct with expiry date.
func (c *Cacher[C, T]) packValue(val T) *value[T] {
	v := value[T]{
		val: val,
	}
	if c.ttl != 0 {
		v.expiry = time.Now().Unix() + c.ttl
	}
	return &v
}

// Delete is used to delete the input key from current
// Cacher instance. It doesn't return anything. If there
// is no such key, Delete is a no-op.
func (c *Cacher[C, T]) Delete(key C) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cacheMap, key)
}

// DeleteSome is used to delete keys which satisfied a
// particular condition determined via SegrigatorFunc.
func (c *Cacher[C, T]) DeleteSome(cond SegrigatorFunc[T]) {
	if cond == nil {
		cond = func(T) bool { return true }
	}
	c.deleteSome(cond)
}

func (c *Cacher[C, T]) deleteSome(cond SegrigatorFunc[T]) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.cacheMap {
		if !cond(v.val) {
			continue
		}
		delete(c.cacheMap, k)
	}
}

// The Reset function deletes the current cache map
// and reallocates an empty one in place of it.
// Use it if you want to delete all keys at once.
// It doesn't return anything.
func (c *Cacher[C, T]) Reset() {
	c.status = cacherReset
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap = make(map[C]*value[T])
}

// NumKeys counts the number of keys present in the
// current Cacher instance and returns that count.
func (c *Cacher[C, T]) NumKeys() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.cacheMap)
}
