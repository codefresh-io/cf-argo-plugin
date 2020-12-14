package bash

import (
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
)

type CommandExecutor struct {
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	output, err := exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId)).Output()
	if err != nil {
		fmt.Printf("Failed to export env id: %s, %s\n", err.Error(), output)
	}
	output, err = exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ACTIVITY_ID=\"%s\"", activity.ActivityId)).Output()
	if err != nil {
		fmt.Printf("Failed to export activity id: %s, result %s\n", err.Error(), output)
	}
}
