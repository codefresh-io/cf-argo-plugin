package sync

import (
	"cf-argo-plugin/pkg/builder"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)


var syncArgs = &builder.SyncArgs{}

var Cmd = &cobra.Command{
	Use:   "sync [app]",
	Short: "Trigger a sync for Argo app",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		b := builder.New()
		err := b.Auth(context.PluginArgoCredentials.Host, context.PluginArgoCredentials.Username, context.PluginArgoCredentials.Password)
		if err != nil {
			return err
		}
		b.ExportExternalUrl(context.PluginArgoCredentials.Host, name)
		b.Sync(syncArgs, name)

		fmt.Println(strings.Join(b.GetLines()[:], "\n"))

		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}

		return nil
	},
}

func init() {
	f := Cmd.Flags()
	f.BoolVar(&syncArgs.Sync, "sync", true, "Specify whether to trigger an ArgoCD sync")
	f.BoolVar(&syncArgs.WaitHealthy, "wait-healthy", false, "Specify whether to wait for sync to be completed (in canary consider wait for suspended status)")
	f.BoolVar(&syncArgs.WaitForSuspend, "wait-suspend", false, "Specify whether to wait for application suspended status")
}
