package processResult

import (
	"cf-argo-plugin/pkg/argo"
	"cf-argo-plugin/pkg/bash"
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var waitRolloutArgsOptions struct {
	PipelineId string
	BuildId    string
}

var WaitRolloutCmd = &cobra.Command{
	Use:   "wait-rollout [app]",
	Short: "Wait rollout",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		fmt.Println("Start execute wait for rollout " + name)

		argoApi := argo.Argo{
			Host:     context.PluginArgoCredentials.Host,
			Username: context.PluginArgoCredentials.Username,
			Password: context.PluginArgoCredentials.Password,
			Token:    context.PluginArgoCredentials.Token,
		}

		cf := codefresh.New(&codefresh.ClientOptions{
			Token: context.PluginCodefreshCredentials.Token,
			Host:  context.PluginCodefreshCredentials.Host,
		})

		historyId, _ := argoApi.GetLatestHistoryId(name)
		fmt.Println(fmt.Sprintf("Current history id %v", historyId))
		start := time.Now()
		for {
			currentHistoryId, _ := argoApi.GetLatestHistoryId(name)
			fmt.Println(fmt.Sprintf("Current in loop history id %v", currentHistoryId))
			// we identify new rollout
			if currentHistoryId > historyId {
				fmt.Println(fmt.Sprintf("Found new history id %v", currentHistoryId))

				// ignore till we will handle it in correct way, 500 code mean that history not found and we shouldnt break pipeline
				_, updatedActivities := cf.SendMetadata(&codefresh.ArgoApplicationMetadata{
					Pipeline:        waitRolloutArgsOptions.PipelineId,
					BuildId:         waitRolloutArgsOptions.BuildId,
					HistoryId:       currentHistoryId,
					ApplicationName: name,
				})

				if updatedActivities != nil {

					err, activity := filterActivity(name, updatedActivities)
					if err == nil {
						bash.CommandExecutor{}.ExportGitopsInfo(activity)
					}
				}

				return nil
			}

			time.Sleep(10 * time.Second)

			elapsed := time.Now().Sub(start)
			if elapsed.Minutes() >= 5 {
				return nil
			}
		}
	},
}

func init() {
	f := WaitRolloutCmd.Flags()
	f.StringVar(&waitRolloutArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&waitRolloutArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
}

func filterActivity(applicationName string, updatedActivities []codefresh.UpdatedActivity) (error, codefresh.UpdatedActivity) {
	var rolloutActivity codefresh.UpdatedActivity
	for _, activity := range updatedActivities {

		if activity.EnvironmentName == applicationName {
			return nil, activity
		}
	}
	return errors.New(fmt.Sprintf("can't find activity with app name %s", applicationName)), rolloutActivity
}
