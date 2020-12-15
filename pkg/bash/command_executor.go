package bash

import (
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
	"strings"
)

type CommandExecutor struct {
}

func execCommand(command string) {
	cmd := exec.Command("/bin/bash", "-c", "echo "+command+" >> /meta/env_vars_to_export")
	err := cmd.Run()
	fmt.Println(strings.Join(cmd.Args, ","))
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_NAME=%s", activity.EnvironmentName))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_ID=%s", activity.EnvironmentId))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ACTIVITY_ID=%s", activity.ActivityId))
}
