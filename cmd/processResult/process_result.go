package processResult

import (
	argo "cf-argo-plugin/pkg/argo"
	"cf-argo-plugin/pkg/builder"
	codefresh "cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

var processResultArgsOptions struct {
	PipelineId 				string
	BuildId    				string
	ExportOutGitopsCommand  string
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
                exportGitopsInfo(activity)
            }
		}

		return nil
	},
}

func exportGitopsInfo(activity codefresh.UpdatedActivity) {
	fmt.Println("Starting export Gitops Info")
	b := builder.New()
	b.ExportGitopsInfo(activity.EnvironmentId, activity.ActivityId)
	resultExportCommands := strings.Join(b.GetExportLines()[:], " ")
	cmd :=exec.Command(resultExportCommands)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error while trying to updaet Gitops Info: \"%q\"", err.Error())
		return
	}
	fmt.Println("Gitops Info exported successfully")
}

func filterActivity(applicationName string, updatedActivities []codefresh.UpdatedActivity) (error, codefresh.UpdatedActivity) {
	fmt.Println("filterActivity")
	var rolloutActivity codefresh.UpdatedActivity
	for _, activity := range updatedActivities {

		if activity.EnvironmentName == applicationName {
			fmt.Println("filterActivity successfully")
			return nil, activity
		}
	}
	return errors.New(fmt.Sprintf("can't find activity with app name %s", applicationName)), rolloutActivity
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&processResultArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&processResultArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
	f.StringVar(&processResultArgsOptions.ExportOutGitopsCommand, "out-export-file", "", "Write export commands to file")
}
