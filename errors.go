package sequencer

import "errors"

var ErrLimitReached error = errors.New("Limit Reached")
var ErrNotImplemented error = errors.New("Implementation pending")
var ErrLockNotGranted error = errors.New("Could not get Lock")
