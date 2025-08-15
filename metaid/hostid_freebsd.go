//go:build freebsd

package metaid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("kern.hostuuid")
}
