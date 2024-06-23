# RunOasis

**RunOasis** is a highly configurable task runner designed to execute tasks in isolated environments. It leverages `bubblewrap`
for sandboxing tasks and `nix-shell` for managing and isolating dependencies for each task, ensuring that every task runs in
a consistent, controlled, and secure setting.

## Key Features

1. **Secure Isolation**:
   - Utilizes `bubblewrap` to create a sandboxed environment, isolating each task from the host system and other tasks.
   - Ensures minimal permissions and controls are granted to each task, reducing the risk of interference or security
     breaches.

2. **Reproducible Environments**:
   - Uses `nix-shell` to provide consistent and reproducible environments for each task.
   - Ensures that the same environment can be recreated reliably, facilitating debugging and development.

3. **Task Flexibility**:
   - Supports a wide range of tasks, from generating cryptocurrency seed phrases to running complex CI/CD pipelines.
   - Easily configurable through simple task definitions, specifying required packages, commands, and environment variables.

## Example Use Case: Secure Seed Phrase Generation

Here's a Proof of Concept (PoC) on how you can define and run a task to generate a BIP39 seed phrase, encrypt it using OpenSSL
with PBKDF2, and store it securely in an isolated environment. This example can be enhanced to send the encrypted seed phrase
to a shared folder on the host machine for further processing or storage.

### Task Definition

```go
package main

import (
    "github.com/aledsdavies/runoasis/core"
)

func main() {
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
```

## Summary

- **RunOasis**: A secure, isolated environment for executing tasks with a focus on security and reproducibility.
- **Key Features**: Secure isolation, reproducible environments, flexible task definitions, and strong security and compliance
  measures.
- **Example Use Case**: Secure generation and encryption of BIP39 seed phrases.
- **Task Runner Implementation**: Ensures tasks are executed securely and efficiently using `nix-shell` and `bubblewrap`.

## Note

This project is a work-in-progress (WIP) and is subject to changes. Contributions and feedback are welcome.
