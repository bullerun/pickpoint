package postgres

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrNothingHappened    = errors.New("the command could not be executed")
	ErrDifferentUserIDs   = errors.New("orders belong to different recipients")
)
