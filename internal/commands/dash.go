package commands

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/twosson/kubeapt/internal/dash"
	"log"
	"os"
	"os/signal"
)

func newDashCmd() *cobra.Command {
	var namespace string
	var uiURL string

	dashCmd := &cobra.Command{
		Use:   "dash",
		Short: "Show dashboard",
		Long:  "Kubeapt Kubernetes dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)

			runCh := make(chan bool, 1)

			go func() {
				if err := dash.Run(ctx, namespace, uiURL); err != nil {
					log.Print(err)
					os.Exit(1)
				}

				runCh <- true
			}()

			select {
			case <-sigCh:
				log.Print("Shutting dashboard down due to interrupt")
			case <-runCh:
				log.Print("Dashboard has exited")
			}
		},
	}

	dashCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	dashCmd.Flags().StringVar(&uiURL, "ui-url", "", "UI URL")

	return dashCmd
}
