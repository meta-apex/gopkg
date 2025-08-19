//go:build !darwin && !linux && !freebsd && !openbsd && !windows

package metaid

import "errors"

func readPlatformMachineID() (string, error) {
	return "", errors.New("not implemented")
}
