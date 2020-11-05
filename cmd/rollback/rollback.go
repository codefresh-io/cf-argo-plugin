package rollback

import (
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rollbackOptions struct {
	Code int
	// if customer want do rollback on failed argocd sync
	NeedRollback bool
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

		if rollbackOptions.Code > 0 && rollbackOptions.NeedRollback {
			log.Println(fmt.Sprintf("Start do rollback because code is \"%v\"", rollbackOptions.Code))

			rollbackResult, _ := cf.RollbackToStable(name, codefresh.Rollback{
				ContextName:     context.PluginCodefreshCredentials.Integration,
				ApplicationName: name,
			})

			if rollbackResult != nil {
				log.Println(fmt.Sprintf("Run rollback process pipeline, build link https://g.codefresh.io/build/%s", rollbackResult.BuildId))
			}

			return errors.New(fmt.Sprintf("ArgoCD app wait fails with the error code \"%v\"", rollbackOptions.Code))
		}

		log.Println("Skip rollback execution")

		return nil
	},
}

func init() {
	f := Cmd.Flags()
	f.IntVar(&rollbackOptions.Code, "code", 0, "Error code from sync execution")
	f.BoolVar(&rollbackOptions.NeedRollback, "needRollback", false, "Execute rollback if sync is failed")
}
