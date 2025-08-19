//go:build openbsd

package metaid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("hw.uuid")
}
