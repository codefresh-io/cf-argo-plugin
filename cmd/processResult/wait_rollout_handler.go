package processResult

import (
	"cf-argo-plugin/pkg/argo"
	"cf-argo-plugin/pkg/bash"
	"cf-argo-plugin/pkg/codefresh"
	"cf-argo-plugin/pkg/context"
	"cf-argo-plugin/pkg/util"
	"errors"
	"fmt"
	"time"
)

const (
	RETRIES_COUNT = 3
	TIMEOUT       = 15
)

type WaitRolloutHandler struct {
	codefresh codefresh.Codefresh
	argo      argo.Argo
}

func GetWaitRolloutHandler() *WaitRolloutHandler {
	return &WaitRolloutHandler{codefresh: codefresh.New(&codefresh.ClientOptions{
		Token: context.PluginCodefreshCredentials.Token,
		Host:  context.PluginCodefreshCredentials.Host,
	}), argo: argo.New(&argo.ClientOptions{
		Host:     context.PluginArgoCredentials.Host,
		Username: context.PluginArgoCredentials.Username,
		Password: context.PluginArgoCredentials.Password,
		Token:    context.PluginArgoCredentials.Token,
	})}
}

func (waitRolloutHandler *WaitRolloutHandler) processNewHistoryId(historyId int64, name string) error {
	fmt.Println(fmt.Sprintf("Found new history id %v", historyId))

	// wait before activity on codefresh will be created
	time.Sleep(15 * time.Second)

	// ignore till we will handle it in correct way, 500 code mean that history not found and we shouldnt break pipeline
	_, updatedActivities := waitRolloutHandler.codefresh.SendMetadata(&codefresh.ArgoApplicationMetadata{
		Pipeline:        waitRolloutArgsOptions.PipelineId,
		BuildId:         waitRolloutArgsOptions.BuildId,
		HistoryId:       historyId,
		ApplicationName: name,
	})

	if updatedActivities != nil {

		err, activity := util.FilterActivity(name, updatedActivities)
		if err == nil {
			fmt.Println(fmt.Sprintf("Start export gitops information"))
			bash.CommandExecutor{}.ExportGitopsInfo(activity)
			return nil
		}
	} else {
		fmt.Println(fmt.Sprintf("Failed to export gitops info, because didnt find activity with history id %v, retrying", historyId))
	}

	return errors.New("failed to export gitops info")
}

func (waitRolloutHandler *WaitRolloutHandler) Handle(name string) error {
	retriesCount := 0

	fmt.Println("Start execute wait for rollout " + name)

	historyId, _ := waitRolloutHandler.argo.GetLatestHistoryId(name)
	start := time.Now()
	for {
		currentHistoryId, _ := waitRolloutHandler.argo.GetLatestHistoryId(name)
		// we identify new rollout
		if currentHistoryId > historyId {
			err := waitRolloutHandler.processNewHistoryId(currentHistoryId, name)
			if err == nil {
				return nil
			}
			// if err exist we need continue to process
			retriesCount++
		}

		time.Sleep(10 * time.Second)

		elapsed := time.Now().Sub(start)
		if elapsed.Minutes() >= TIMEOUT || retriesCount >= RETRIES_COUNT {
			fmt.Println("Stop wait for rollout because retries time exceed")
			return nil
		}
	}
}
