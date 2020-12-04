package rollout

import (
	"cf-argo-plugin/pkg/builder"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rolloutArgs = &builder.RolloutArgs{}

var Cmd = &cobra.Command{
	Use:   "rollout [app]",
	Short: "Promote Argo rollout",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		b := builder.New()

		// we are using healthy validation but without creds
		if context.PluginCodefreshCredentials.Token == "" && context.PluginCodefreshCredentials.Host == "" && rolloutArgs.WaitHealthy {
			return errors.New("For use wait_healthy flag you should provide context or argo credentials argument")
		}

		if context.PluginArgoCredentials.Token == "" && context.PluginCodefreshCredentials.Host != "" {
			err := b.Auth(context.PluginArgoCredentials.Host, context.PluginArgoCredentials.Username, context.PluginArgoCredentials.Password)
			if err != nil {
				return err
			}
		}
		b.ExportExternalUrl(context.PluginArgoCredentials.Host, name)
		b.Rollout(rolloutArgs, name, context.PluginArgoCredentials.Token, context.PluginArgoCredentials.Host)

		resultCommands := strings.Join(b.GetLines()[:], "\n")
		resultExportCommands := strings.Join(b.GetExportLines()[:], "\n")

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
	f.StringVar(&rolloutArgs.KubernetesContext, "k8s-context", "", "The k8s context name as it show in the k8s integration.")
	f.StringVar(&rolloutArgs.RolloutName, "rollout-name", "", "The name of the rollout to be promoted.")
	f.StringVar(&rolloutArgs.RolloutNamespace, "rollout-namespace", "default", "The namespace of the rollout to be promoted.")
	f.BoolVar(&rolloutArgs.WaitHealthy, "wait-healthy", true, "Specify whether to wait for sync to be completed (in canary consider wait for suspended status)")

	_ = cobra.MarkFlagRequired(f, "k8s-context")
	_ = cobra.MarkFlagRequired(f, "rollout-name")
}
