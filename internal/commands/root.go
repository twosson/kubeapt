package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	// remove timestamp from log
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// Execute executes apt.
func Execute(gitCommit string, buildTime string) {
	rootCmd := newRoot(gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRoot(gitCommit string, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "apt",
		Short: "apt is the Kubeapt CLI",
	}

	rootCmd.AddCommand(newDashCmd())
	rootCmd.AddCommand(newVersionCmd(gitCommit, buildTime))

	return rootCmd
}
