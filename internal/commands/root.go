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
func Execute() {
	rootCmd := newRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "apt",
		Short: "apt is the Kubeapt CLI",
	}

	rootCmd.AddCommand(newDashCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
