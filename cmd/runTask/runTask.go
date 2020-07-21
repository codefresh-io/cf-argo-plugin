package runTask

import (
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "run-task [env name]",
	Short: "Run codefresh task for sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		codefreshApi := codefresh.New(&codefresh.ClientOptions{
			Token: context.PluginCodefreshCredentials.Token,
			Host:  context.PluginCodefreshCredentials.Host,
		})

		result, err := codefreshApi.StartSyncTask(name)

		if err != nil {
			return err
		}

		fmt.Printf("Build id: %s", result.BuildId)

		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a env name argument")
		}

		return nil
	},
}
