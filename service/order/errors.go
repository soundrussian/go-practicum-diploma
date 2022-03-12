package order

import "errors"

var (
	ErrOrderInvalid    = errors.New("order id is invalid")
	ErrConflict        = errors.New("order id was already uploaded by another user")
	ErrAlreadyAccepted = errors.New("order id has been uploaded by this user earlier")
)
