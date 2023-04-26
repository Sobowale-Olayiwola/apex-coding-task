package db

import "errors"

var (
	ErrNoRecord         = errors.New("no matching record found")
	ErrInsufficientFund = errors.New("user has insufficient amount. kindly fund your wallet")
	ErrActiveSession    = errors.New("user has an active session")
)
