package bash

import (
	"bytes"
	"cf-argo-plugin/pkg/codefresh"
	"fmt"
	"os/exec"
)

type CommandExecutor struct {
}

func execCommand(command string) {
	cmd := exec.Command("/bin/bash", "cf_export", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
	fmt.Printf("Command result: %q\n", out.String())
}

func execCommand2(command string) {
	cmd := exec.Command("/bin/bash", "cf_export", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
	fmt.Printf("Command result: %q\n", out.String())
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ACTIVITY_ID=%s", activity.ActivityId))
}
