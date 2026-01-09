package core

import "errors"

var ErrNotFound = errors.New("element not found")
var ErrFailedAdding = errors.New("failed to add")
var ErrFailedUpdate = errors.New("failed to update element")
var ErrFailedDelete = errors.New("failed to delete element")
