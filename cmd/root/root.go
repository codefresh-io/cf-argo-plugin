package root

import (
	"cf-argo-plugin/cmd/processResult"
	"cf-argo-plugin/cmd/rollback"
	"cf-argo-plugin/cmd/rollout"
	"cf-argo-plugin/cmd/runTask"
	"cf-argo-plugin/cmd/sync"
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type authContext struct {
	CodefreshToken       string
	CodefreshHost        string
	CodefreshIntegration string

	ArgoUsername string
	ArgoPassword string
	ArgoHost     string
	ArgoToken    string
}

var pluginAuthContext = &authContext{}

var rootCmd = &cobra.Command{
	Use:   "cf-argo-plugin",
	Short: "Codefresh plugin for easy interact with argocd",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}

		return nil
	},
	PersistentPreRunE: fetchArgoCredentials,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&pluginAuthContext.CodefreshToken, "cf-token", "", "Api token from Codefresh")
	pf.StringVar(&pluginAuthContext.CodefreshHost, "cf-host", "https://g.codefresh.io", "Host of Codefresh")
	pf.StringVar(&pluginAuthContext.CodefreshIntegration, "cf-integration", "", "Name of Argo integration on Codefresh")

	pf.StringVar(&pluginAuthContext.ArgoUsername, "argo-username", "", "Username for argo cd, use only if you not provide integration")
	pf.StringVar(&pluginAuthContext.ArgoPassword, "argo-password", "", "Password for argo cd, use only if you not provide integration")
	pf.StringVar(&pluginAuthContext.ArgoToken, "argo-token", "", "Token for argo cd, use only if you not provide integration")
	pf.StringVar(&pluginAuthContext.ArgoHost, "argo-host", "", "Host for argo cd, use only if you not provide integration")

	pf.StringVar(&context.PluginOutConfig.CommandsFile, "out-commands-file", "", "Write main commands to file")
	pf.StringVar(&context.PluginOutConfig.ExportOutUrlCommand, "out-export-file", "", "Write export commands to file")

	rootCmd.AddCommand(sync.Cmd)
	rootCmd.AddCommand(rollout.Cmd)
	rootCmd.AddCommand(runTask.Cmd)
	rootCmd.AddCommand(processResult.Cmd)
	rootCmd.AddCommand(rollback.Cmd)

}

func fetchArgoCredentials(cmd *cobra.Command, args []string) error {
	context.PluginCodefreshCredentials.Host = pluginAuthContext.CodefreshHost
	context.PluginCodefreshCredentials.Token = pluginAuthContext.CodefreshToken

	if pluginAuthContext.ArgoUsername != "" && pluginAuthContext.ArgoPassword != "" && pluginAuthContext.ArgoHost != "" {
		context.PluginArgoCredentials.Host = pluginAuthContext.ArgoHost
		context.PluginArgoCredentials.Username = pluginAuthContext.ArgoUsername
		context.PluginArgoCredentials.Password = pluginAuthContext.ArgoPassword
	} else if pluginAuthContext.CodefreshToken != "" && pluginAuthContext.CodefreshIntegration != "" {
		codefreshApi := codefresh.New(&codefresh.ClientOptions{
			Token: pluginAuthContext.CodefreshToken,
			Host:  pluginAuthContext.CodefreshHost,
		})
		integration, err := codefreshApi.GetIntegration(pluginAuthContext.CodefreshIntegration)
		if err != nil {
			return fmt.Errorf("failed to retrive argo integration, %s", err.Error())
		}

		context.PluginArgoCredentials.Host = integration.Data.Url
		context.PluginArgoCredentials.Username = integration.Data.Username
		context.PluginArgoCredentials.Password = integration.Data.Password
		context.PluginArgoCredentials.Token = integration.Data.Token
	} else {
		return fmt.Errorf("can`t resolve argo context, you should provide correct codefresh or raw argo context")
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
