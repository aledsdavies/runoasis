package core

import (
	"fmt"
	"strings"
)

func buildNixTask(task Task) string {
	var environentVars string
	for key, value := range task.Variables {
		environentVars += fmt.Sprintf("%s = \"%s\";\n", key, value)
	}

	// Define a basic .nix file content
	nixContent := fmt.Sprintf(`
{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [ %s ];

  %s

  shellHook = ''
    %s
  '';
}`, strings.Join(task.Packages, " "), environentVars, strings.Join(task.Commands, "\n"))

	return nixContent
}
