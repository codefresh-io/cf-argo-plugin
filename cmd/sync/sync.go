package sync

import (
	"cf-argo-plugin/pkg/builder"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var syncArgs = &builder.SyncArgs{}

var Cmd = &cobra.Command{
	Use:   "sync [app]",
	Short: "Trigger a sync for Argo app",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		b := builder.New()
		if context.PluginArgoCredentials.Token == "" {
			err := b.Auth(context.PluginArgoCredentials.Host, context.PluginArgoCredentials.Username, context.PluginArgoCredentials.Password)
			if err != nil {
				return err
			}
		}

		b.ExportExternalUrl(context.PluginArgoCredentials.Host, name)
		b.Sync(syncArgs, name, context.PluginArgoCredentials.Token, context.PluginArgoCredentials.Host, context.PluginCodefreshCredentials.Integration)

		resultCommands := strings.Join(b.GetLines()[:], "\n")
		resultExportCommands := strings.Join(b.GetExportLines()[:], "\n")

		if syncArgs.Debug {
			fmt.Println("Command to execute: " + resultCommands)
		}

		if context.PluginOutConfig.CommandsFile != "" {
			file, err := os.Create(context.PluginOutConfig.CommandsFile)
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = file.WriteString(resultCommands)

			if err != nil {
				return err
			}
		} else {
			fmt.Println(resultCommands)
		}

		if context.PluginOutConfig.ExportOutUrlCommand != "" {
			file, err := os.Create(context.PluginOutConfig.ExportOutUrlCommand)
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = file.WriteString(resultExportCommands)

			if err != nil {
				return err
			}
		}

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
	f.BoolVar(&syncArgs.Prune, "prune", false, "Allow deleting unexpected resources")
	f.BoolVar(&syncArgs.WaitForSuspend, "wait-suspend", false, "Specify whether to wait for application suspended status")
	f.BoolVar(&syncArgs.Debug, "debug", false, "Debug argocd command ( print commands to output )")
	f.StringVar(&syncArgs.AdditionalFlags, "additional-flags", "", "Specify additional flags , like --grpc-web , so on")
	f.StringVar(&syncArgs.WaitAdditionalFlags, "wait-additional-flags", "", "Specify additional flags for wait command, like --timeout , so on")
	f.StringVar(&syncArgs.Revision, "revision", "", "Sync to a specific revision. Preserves parameter overrides")

}
