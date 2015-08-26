package cbes

import (
    "log"
    "io"
)

type Log struct {
    *log.Logger
}

// set io.Writer to create a Logger.
func NewLog (out io.Writer) *Log {
    _log := new(Log)
    _log.Logger = log.New(out, "[CBES]", 1e9)

    return _log
}