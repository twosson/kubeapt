package main

import "github.com/twosson/kubeapt/internal/commands"

var (
	gitCommit = "(unknown-commit)"
	buildTime = "(unknown-buildtime)"
)

func main() {
	commands.GitCommit = gitCommit
	commands.BuildTime = buildTime

	commands.Execute()
}
