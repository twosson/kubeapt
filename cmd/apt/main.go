package main

import "github.com/twosson/kubeapt/internal/commands"

var (
	gitCommit = "(dev-commit)"
	buildTime = "(dev-buildtime)"
)

func main() {
	commands.Execute(gitCommit, buildTime)
}
