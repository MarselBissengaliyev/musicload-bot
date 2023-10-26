package downloader

import "errors"

var (
	ErrDurationLimitExceeded = errors.New("Request Entity Too Large")
)