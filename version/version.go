package version

import (
	"runtime"
)

// Version is the current version of the application

const version = "0.11.0"

func GetSDKVersion() string {
	return version
}

func GetGoVersion() string {
	return runtime.Version()
}

func GetOS() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}
