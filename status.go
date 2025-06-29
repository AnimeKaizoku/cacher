package cacher

type status int

const (
	noop status = iota
	cleanerBusy
	cacherReset
	cacherDeleted
)

type CleaningMode int

const (
	CleaningNone CleaningMode = iota
	CleaningCentral
	CleaningLocal
)
