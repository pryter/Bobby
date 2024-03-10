package db

import "errors"

var (
	ErrInvalidMacAddr = errors.New("unable to parse mac_addr")
)
