package api

import (
	"errors"
)

// ErrorNotFound - 404 Not Found
var ErrorNotFound = errors.New("Not Found")

// ErrorInvalidContentType -
var ErrorInvalidContentType = errors.New("Invalid ContentType")

// ErrorInvalidOperation is operation without "create", "update" or "delete"
var ErrorInvalidOperation = errors.New("Invalid Operation")

// ErrorSync is error during the sync operation
var ErrorSync = errors.New("Sync Error")

// ErrorSequence is error during the sync operation
var ErrorSequence = errors.New("Sequence Error")

// ErrorOutOfSync events
var ErrorOutOfSync = errors.New("Out of sync events")
