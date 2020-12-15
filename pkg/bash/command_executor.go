package bash

import (
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
)

type CommandExecutor struct {
}

func execCommand(command string) {
	cmd := exec.Command("bash", "-c", "echo "+command+" >> /meta/env_vars_to_export")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_NAME=\"%s\"", activity.EnvironmentName))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ACTIVITY_ID=%s", activity.ActivityId))
}
