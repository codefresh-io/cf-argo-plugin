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

type (
	WaitRollout interface {
		Wait(name string, pipelineId string, buildId string) error
	}

	waitRollout struct {
	}
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
			Pipeline:        waitRolloutArgsOptions.PipelineId,
			BuildId:         waitRolloutArgsOptions.BuildId,
			HistoryId:       historyId,
			ApplicationName: name,
		})

		if updatedActivities != nil {

			err, activity := filterActivity(name, updatedActivities)
			if err == nil {
				bash.CommandExecutor{}.ExportGitopsInfo(activity)
			}
		}

		return nil
	},
}

func init() {
	f := WaitRolloutCmd.Flags()
	f.StringVar(&waitRolloutArgsOptions.PipelineId, "pipeline-id", "", "Pipeline id where argo sync was executed")
	f.StringVar(&waitRolloutArgsOptions.BuildId, "build-id", "", "Build id where argo sync was executed")
}

func (wr *waitRollout) Wait(name string, pipelineId string, buildId string) error {
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
				Pipeline:        pipelineId,
				BuildId:         buildId,
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