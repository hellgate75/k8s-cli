package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func ConfigDir() string {
	var home = UserHomeDir()
	if home == "" {
		home = fmt.Sprintf(".%c", os.PathSeparator)
	}
	home = fmt.Sprintf("%s%c%s", home, os.PathSeparator, ".k8s-cli")
	return home
}

func FixOutputType(oType string) string {
	if strings.ToLower(oType) == "yaml" {
		return strings.ToLower(oType)
	}
	return "json"
}