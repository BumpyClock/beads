package main

import (
	"context"
	"errors"
)

// DoltServerHandle is a stub; Dolt backend has been removed.
type DoltServerHandle struct{}

// DoltDefaultSQLPort is the default SQL port for dolt server
const DoltDefaultSQLPort = 3306

// DoltDefaultRemotesAPIPort is the default remotesapi port for dolt server
const DoltDefaultRemotesAPIPort = 50051

// ErrDoltRequiresCGO is returned when dolt features are requested.
var ErrDoltRequiresCGO = errors.New("dolt backend has been removed; see github.com/BumpyClock/beads-dolt")

// StartDoltServer returns an error; Dolt backend has been removed.
func StartDoltServer(ctx context.Context, dataDir, logFile string, sqlPort, remotePort int) (*DoltServerHandle, error) {
	return nil, ErrDoltRequiresCGO
}

// Stop is a no-op stub.
func (h *DoltServerHandle) Stop() error {
	return nil
}

// SQLPort returns 0; Dolt backend has been removed.
func (h *DoltServerHandle) SQLPort() int {
	return 0
}

// RemotesAPIPort returns 0; Dolt backend has been removed.
func (h *DoltServerHandle) RemotesAPIPort() int {
	return 0
}

// Host returns empty string; Dolt backend has been removed.
func (h *DoltServerHandle) Host() string {
	return ""
}

// DoltServerAvailable returns false; Dolt backend has been removed.
func DoltServerAvailable() bool {
	return false
}


