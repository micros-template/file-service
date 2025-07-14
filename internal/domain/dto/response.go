package dto

import "errors"

var (
	Err_INTERNAL_SAVE_PROFILE_IMAGE   = errors.New("failed to save profile image")
	Err_INTERNAL_REMOVE_PROFILE_IMAGE = errors.New("failed to remove profile image")
)
