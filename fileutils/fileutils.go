// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /fileutils.go
// Original timestamp: 2024/06/22 13:24

package fileutils

import (
	"os"
	"runtime"
	"strings"

	cerr "github.com/jeanfrancoisgratton/customError/v2"
)

// Check which type of filesystem the mountpoint is.
// Currently only supports MacOS and linux; I've broken my Windows test VM and have no time for it
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

	return "", &cerr.CustomError{Fatality: cerr.Continuable, Title: "WIP", Message: "Unsupported for now.. stay tuned"}

	//var stat unix.Statfs_t
	//err := unix.Statfs(mountpoint, &stat)
	//if err != nil {
	//	return "", &cerr.CustomError{Title: "Error with Statfs syscall", Message: err.Error()}
	//}
	//
	//// The Type field indicates the file system type
	//fsType := stat.Type
	//
	//// Mapping file system type to a human-readable string
	//fsTypeName := ""
	//switch fsType {
	//case unix.MNT_ASYNC:
	//	fsTypeName = "async"
	//case unix.MNT_LOCAL:
	//	fsTypeName = "local"
	//// Add more cases as needed for different file system types
	//default:
	//	fsTypeName = fmt.Sprintf("unknown (%d)", fsType)
	//}
	//
	//return fsTypeName, nil
}
