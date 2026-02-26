package errs

import "errors"

// ErrClosedConn occurs when the connection already closed.
var ErrClosedConn = errors.New("closed connection")

// ErrInvalidSyntax occurs when the syntax of something is wrong.
var ErrInvalidSyntax = errors.New("invalid syntax")
