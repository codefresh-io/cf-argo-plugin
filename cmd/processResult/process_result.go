package processResult

import (
	"cf-argo-plugin/pkg/argo"
	"cf-argo-plugin/pkg/bash"
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"github.com/spf13/cobra"
)

var processResultArgsOptions struct {
	PipelineId             string
	BuildId                string
	ExportOutGitopsCommand string
}

var Cmd = &cobra.Command{
	Use:   "process-result [app]",
	Short: "Process plugin execution result",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		argoApi := argo.Argo{
			Host:     context.PluginArgoCredentials.Host,
			Username: context.PluginArgoCredentials.Username,
			Password: context.PluginArgoCredentials.Password,
			Token:    context.PluginArgoCredentials.Token,
		}

		historyId, _ := argoApi.GetLatestHistoryId(name)

		cf := codefresh.New(&codefresh.ClientOptions{
			Token: context.PluginCodefreshCredentials.Token,
			Host:  context.PluginCodefreshCredentials.Host,
		})

		// ignore till we will handle it in correct way, 500 code mean that history not found and we shouldnt break pipeline
		_, updatedActivities := cf.SendMetadata(&codefresh.ArgoApplicationMetadata{
			Pipeline:        processResultArgsOptions.PipelineId,
			BuildId:         processResultArgsOptions.BuildId,
			HistoryId:       historyId,
			ApplicationName: name,
		})

		if updatedActivities != nil && processResultArgsOptions.ExportOutGitopsCommand != "" {

			err, activity := filterActivity(name, updatedActivities)
			if err == nil {
				bash.CommandExecutor{}.ExportGitopsInfo(activity)
			}
		}

		return nil
	},
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&processResultArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&processResultArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
	f.StringVar(&processResultArgsOptions.ExportOutGitopsCommand, "out-export-file", "", "Write export commands to file")
}
