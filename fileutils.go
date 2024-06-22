// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /fileutils.go
// Original timestamp: 2024/06/22 13:24

package helperFunctions

import (
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"os"
	"runtime"
	"strings"
	"syscall"
)

// Check which type of filesystem the mountpoint is.
// Currently only supports MacOS and linux
func GetFStype(mountpoint string) (string, *cerr.CustomError) {
	switch os := runtime.GOOS; os {
	case "darwin":
		return getMacOSMountPointType(mountpoint)
	case "linux":
		return getLinuxMountPointType(mountpoint)
	default:
		return "", &cerr.CustomError{Fatality: cerr.Continuable,
			Title: "Unsupported operating system", Message: "os type: " + os,
		}
	}
}

// Gets the type of the filesystem by reading the /proc/mounts
func getLinuxMountPointType(mountpoint string) (string, *cerr.CustomError) {
	if runtime.GOOS != "linux" {
		return "", &cerr.CustomError{Fatality: cerr.Continuable, Title: "OS not supported"}
	}

	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return "", &cerr.CustomError{Title: "Error reading /proc/mounts"}
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if fields[1] == mountpoint {
			return fields[2], nil
		}
	}

	return "", &cerr.CustomError{Fatality: cerr.Continuable, Title: "mountpoint not found"}
}

func getMacOSMountPointType(mountpoint string) (string, *cerr.CustomError) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(mountpoint, &stat)
	if err != nil {
		return "", &cerr.CustomError{Title: "Error with Stats syscall :", Message: err.Error()}
	}

	return fmt.Sprintf("%s", stat.Fstypename), nil
}
