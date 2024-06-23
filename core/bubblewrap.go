package core

import (
	"os/exec"

	"github.com/alessio/shellescape"
)


func buildBubbleWrapCommand(baseDir, resolvConfPath, shellCmd string, mountPaths []string) *exec.Cmd {
	args := []string{
		"--dev", "/dev", // Create minimal /dev directory
		"--proc", "/proc", // Mount a new proc filesystem
		"--ro-bind", resolvConfPath, "/etc/resolv.conf", // Bind the temporary resolv.conf with Google DNS
		"--bind", baseDir, baseDir, // Bind the base project directory
		"--unshare-all",
		"--share-net",
		"--die-with-parent",
	}

	// Add the --ro-bind arguments for each mount path
	for _, path := range mountPaths {
		safePath := shellescape.Quote(path)
		args = append(args, "--ro-bind", safePath, safePath)
	}

	args = append(args, "--", "sh", "-c", shellCmd)

	return exec.Command("bwrap", args...)
}
