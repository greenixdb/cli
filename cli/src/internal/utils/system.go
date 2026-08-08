package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

type SystemInfo struct {
	OS       string
	Arch     string
	Hostname string
}

func GetSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()
	return SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
	}
}

func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		return cmd.Run() == nil
	}
	return os.Geteuid() == 0
}

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("COMSPEC")
	}
	return shell
}

