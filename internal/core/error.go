package core

import "errors"

var (
	ErrUsernameExists  = errors.New("username_exists")
	ErrStarIDNotExists = errors.New("star_id_not_exists")
)
