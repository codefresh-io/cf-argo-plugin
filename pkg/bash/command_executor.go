package bash

import (
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
)

type CommandExecutor struct {
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId))
	exec.Command("/bin/sh", fmt.Sprintf("cf_export sendArgoMetadata_CF_ACTIVITY_ID=\"%s\"", activity.ActivityId))
}
