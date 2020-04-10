package main

import "github.com/twosson/kubeapt/internal/commands"

var (
	version   = "(dev-version)"
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func main() {
	commands.Execute(version, gitCommit, buildTime)
}
