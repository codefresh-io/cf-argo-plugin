package bash

import (
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
)

type CommandExecutor struct {
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	err := exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId)).Run()
	if err != nil {
		fmt.Printf("Failed to export env id: %s\n", err.Error())
	}
	err = exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ACTIVITY_ID=\"%s\"", activity.ActivityId)).Run()
	if err != nil {
		fmt.Printf("Failed to export activity id: %s\n", err.Error())
	}
}
