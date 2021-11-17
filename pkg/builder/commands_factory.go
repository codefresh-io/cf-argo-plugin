package builder

import "fmt"

type CommandsFactory struct {
}

func GetCommandsFactory() *CommandsFactory {
	return &CommandsFactory{}
}

func (cf *CommandsFactory) CreateWaitRolloutCMD(name string, integration string, skip bool) string {
	skipFlag := ""
	if skip {
		skipFlag = "--skip"
	}

	return fmt.Sprintf(`
		cf-argo-plugin wait-rollout %s %s --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=%s --pipeline-id="$CF_PIPELINE_NAME" --build-id=$CF_BUILD_ID &
        sleep 5s
		`, name, skipFlag, integration)
}
