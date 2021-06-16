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

		retriesCount := 0

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
		start := time.Now()
		for {
			currentHistoryId, _ := argoApi.GetLatestHistoryId(name)
			// we identify new rollout
			if currentHistoryId > historyId {
				fmt.Println(fmt.Sprintf("Found new history id %v", currentHistoryId))

				// wait before activity on codefresh will be created
				time.Sleep(15 * time.Second)

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
						fmt.Println(fmt.Sprintf("Start export gitops information"))
						bash.CommandExecutor{}.ExportGitopsInfo(activity)
						return nil
					}
				} else {
					fmt.Println(fmt.Sprintf("Failed to export gitops info, because didnt find activity with history id %v, retrying", currentHistoryId))
				}
			}

			time.Sleep(10 * time.Second)

			elapsed := time.Now().Sub(start)
			if elapsed.Minutes() >= 15 || retriesCount >= 5 {
				fmt.Println("Stop wait for rollout because retries time exceed")
				return nil
			}

			retriesCount++
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

		if activity.ApplicationName == applicationName {
			return nil, activity
		}
	}
	return errors.New(fmt.Sprintf("can't find activity with app name %s", applicationName)), rolloutActivity
}
