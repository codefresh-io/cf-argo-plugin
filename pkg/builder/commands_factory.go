package builder

import "fmt"

type CommandsFactory struct {
}

func GetCommandsFactory() *CommandsFactory {
	return &CommandsFactory{}
}

func (cf *CommandsFactory) CreateWaitRolloutCMD(name string, integration string) string {
	return fmt.Sprintf(`
		cf-argo-plugin wait-rollout %s --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=%s --pipeline-id="$CF_PIPELINE_NAME" --build-id=$CF_BUILD_ID &
        sleep 5s
		`, name, integration)
}
