//go:build freebsd

package guid

import "syscall"

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("kern.hostuuid")
}
