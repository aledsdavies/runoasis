package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Runner struct {
	mountPaths []string
}

func NewRunner(mountPaths []string) *Runner {
	return &Runner{
		mountPaths: mountPaths,
	}
}

func (r Runner) RunTask(task Task) {
	// Create a temporary directory for the chroot environment
	baseDir, err := os.MkdirTemp("", "task-")
	if err != nil {
		log.Fatalf("[%s] Failed to create temp dir: %v", task.ID, err)
	}
	defer os.RemoveAll(baseDir)

	// ------------------------------------------------------------------------
	// Nix Setup
	// ------------------------------------------------------------------------

	nixContent := buildNixTask(task)

	// Create a temporary .nix file
	nixFile, err := os.CreateTemp(baseDir, "shell-*.nix")
	if err != nil {
		log.Fatalf("[%s] Failed to create temp nix file: %v", task.ID, err)
	}
	defer os.Remove(nixFile.Name())

	if _, err := nixFile.Write([]byte(nixContent)); err != nil {
		log.Fatalf("[%s] Failed to write to temp nix file: %v", task.ID, err)
	}
	nixFile.Close()

	shellCmd := fmt.Sprintf("nix-shell %s --pure", nixFile.Name())

	// ------------------------------------------------------------------------
	// Bubblewrap Setup
	// ------------------------------------------------------------------------

	// Create a temporary resolv.conf file with Google DNS
	resolvConfContent := "nameserver 8.8.8.8\nnameserver 8.8.4.4\n"
	tempResolvConf, err := os.CreateTemp("", "resolv.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary resolv.conf: %v\n", err)
		return
	}
	defer os.Remove(tempResolvConf.Name())

	if _, err := tempResolvConf.Write([]byte(resolvConfContent)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to temporary resolv.conf: %v\n", err)
		return
	}
	tempResolvConf.Close()

	sandboxCommand := buildBubbleWrapCommand(baseDir, tempResolvConf.Name(), shellCmd, r.mountPaths)

	// ------------------------------------------------------------------------
	// Capture outputs of the bubblewrap command
	// ------------------------------------------------------------------------

	stdout, err := sandboxCommand.StdoutPipe()
	if err != nil {
		log.Fatalf("[%s] Failed to get stdout pipe: %v", task.ID, err)
	}

	stderr, err := sandboxCommand.StderrPipe()
	if err != nil {
		log.Fatalf("[%s] Failed to get stderr pipe: %v", task.ID, err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go streamOutput(task.ID, stdout, &wg)
	go streamOutput(task.ID, stderr, &wg)

	if err := sandboxCommand.Start(); err != nil {
		log.Fatalf("[%s] Failed to start command: %v", task.ID, err)
	}

	wg.Wait()

	if err := sandboxCommand.Wait(); err != nil {
		log.Printf("[%s] Command finished with error: %v", task.ID, err)
	} else {
		log.Printf("[%s] Command finished successfully", task.ID)
	}
}

func streamOutput(taskID string, pipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", taskID, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[%s] Error reading output: %v", taskID, err)
	}
}
