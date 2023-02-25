package cacher

type status int

const (
	noop status = iota
	cleanerBusy
	cacherReset
	cacherDeleted
)
