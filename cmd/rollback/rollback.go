package rollback

import (
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rollbackOptions struct {
	Message string
}

var Cmd = &cobra.Command{
	Use:   "rollback [app]",
	Short: "Handle rollback case",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cf := codefresh.New(&codefresh.ClientOptions{
			Token: context.PluginCodefreshCredentials.Token,
			Host:  context.PluginCodefreshCredentials.Host,
		})

		if rollbackOptions.Message != "" {
			rollbackResult, err := cf.RollbackToStable(name, codefresh.Rollback{
				ContextName:     context.PluginCodefreshCredentials.Integration,
				ApplicationName: name,
			})

			if rollbackResult != nil {
				log.Println(fmt.Sprintf("Run rollback process, build link https://g.codefresh.io/build/%s", rollbackResult.BuildId))
			}

			return err
		}

		return nil
	},
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&rollbackOptions.Message, "message", "", "Error message from sync execution")
}
