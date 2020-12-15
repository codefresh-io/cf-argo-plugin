package main

import (
	"os/exec"
	"strings"
)

func main() {
	//root.Execute()

	cmd := exec.Command("/bin/bash", "-c", "echo sendArgoMetadata_CF_ENVIRONMENT_ID=3 >> /tmp/env_vars_to_export")

	err := cmd.Run()

	print(cmd.Path)
	print(strings.Join(cmd.Args, ","))

	if err != nil {
		println(err.Error())
	}

}
