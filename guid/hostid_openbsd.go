//go:build openbsd

package guid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("hw.uuid")
}
