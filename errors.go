package cvsgeolookup

import "errors"

var (
	ErrNotInitialized        = errors.New("Engine is not initialized")
	ErrReadInterfaceRequired = errors.New("Reader interface is required")
	ErrNoBeginField          = errors.New("Cannot find begin field in header")
	ErrNoEndField            = errors.New("Cannot find end field in header")
	ErrNoLantitudeField      = errors.New("Cannot find lantitude field in header")
	ErrNoLongtitudeField     = errors.New("Cannot find longtitude field in header")
	ErrWrongIPFormat         = errors.New("Wrong IP address format")
	ErrIncorrectSegment      = errors.New("Incorrect segment values in data")
	ErrNotFound              = errors.New("Geolocation not found")
)
