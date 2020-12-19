package processResult

import (
	"cf-argo-plugin/pkg/argo"
	"cf-argo-plugin/pkg/bash"
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"github.com/spf13/cobra"
	"time"
)

var WaitRolloutCMD = &cobra.Command{
	Use:   "watch-rollout [app]",
	Short: "Watch for new rollouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

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

		envs, _ := cf.GetEnvironments()

		var existingEnv *codefresh.CFEnvironment

		for _, env := range envs {
			if env.Spec.Application == name {
				existingEnv = &env
			}
		}

		if existingEnv == nil {
			return nil
		}

		historyId, _ := argoApi.GetLatestHistoryId(name)
		start := time.Now()
		for {
			currentHistoryId, _ := argoApi.GetLatestHistoryId(name)
			// we identify new rollout
			if currentHistoryId > historyId {
				// ignore till we will handle it in correct way, 500 code mean that history not found and we shouldnt break pipeline
				_, updatedActivities := cf.SendMetadata(&codefresh.ArgoApplicationMetadata{
					Pipeline:        processResultArgsOptions.PipelineId,
					BuildId:         processResultArgsOptions.BuildId,
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
			if elapsed.Minutes() >= 15 {
				return nil
			}
		}
	},
}

func init() {
	f := WaitRolloutCMD.Flags()
	f.StringVar(&processResultArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&processResultArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
}
