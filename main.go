package main

import (
	"sync"

	"github.com/asdavies/runoasis/core"
)

// TODO: Move this into test
// TODO: Make these part of the external configuration
// TODO: Investigate the setup for different systems because this is limited to NixOS
var mountPaths = []string{
	"/run/current-system/sw/bin",
	"/nix/store",
	"/nix/var/nix",
	"/etc/nix",
	"/etc/static/nix",
	"/etc/ssl",
	"/etc/static/ssl",
}

func main() {
	runner := core.NewRunner(mountPaths)

	// Define an array of jobs
	jobs := []core.Task{
		{ID: "worker_1", Commands: []string{"java --version"}, Packages: []string{"jdk"}},
		{ID: "worker_2", Commands: []string{"go version"}, Packages: []string{"go"}},
		{ID: "worker_3", Commands: []string{"node --version", "go verson"}, Packages: []string{"nodejs_22"}},
		{ID: "worker_4", Commands: []string{"echo \"Hello, ''${NAME}!\""}, Packages: []string{"nodejs_22"}, Variables: map[string]string{"NAME": "Aled"}},
	}

	var wg sync.WaitGroup

	// Run each job concurrently
	for _, job := range jobs {
		wg.Add(1)
		go func(job core.Task) {
			defer wg.Done()
			runner.RunTask(job)
		}(job)
	}

	wg.Wait()

	// WARNING: This is a quick PoC and you should consider security implications for production use.
	task := core.Task{
		ID: "generate-seed",
		Commands: []string{
			"mkdir ./safe",
			// Generate a BIP39 seed phrase using python-mnemonic
			"SEED_PHRASE=$(python3 -c 'from mnemonic import Mnemonic; print(Mnemonic(\"english\").generate(128))')",
			// Encrypt the seed phrase with OpenSSL using PBKDF2
			"echo $SEED_PHRASE | openssl enc -aes-256-cbc -pbkdf2 -salt -out encrypted_seed.txt -pass pass:yourpassword",
			// Copy the encrypted file to a safe location
			"cp encrypted_seed.txt ./safe/encrypted_seed.cpy.txt",
			// Display the encrypted file content (for PoC purposes)
			"cat ./safe/encrypted_seed.cpy.txt",
		},
		Packages: []string{"python312Packages.mnemonic", "openssl"},
	}

	runner.RunTask(task)
}
