package utils

import (
	"runtime"
)

func DiscriminationOS() string {
	return runtime.GOOS
}
