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
	cmd := exec.Command("/bin/bash", "echo 1")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
	fmt.Printf("Command result: %q\n", out.String())
}

func execCommand3(command string) {
	cmd := exec.Command("/bin/echo", "MY_PLUGIN_VAR=SAMPLE_VALUE >> /meta/env_vars_to_export")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute export command: %v\n", err)
	}
}

func (commandExecutor CommandExecutor) ExportGitopsInfo(activity codefresh.UpdatedActivity) {
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ENVIRONMENT_ID=\"%s\"", activity.EnvironmentId))
	execCommand(fmt.Sprintf("sendArgoMetadata_CF_ACTIVITY_ID=%s", activity.ActivityId))
	execCommand2("")
	execCommand3("")
}
