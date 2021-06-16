package builder

import (
	"testing"
)

func TestCreateWaitRolloutCMD(t *testing.T) {
	expectedCMD := `
		cf-argo-plugin wait-rollout test --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=context --pipeline-id=$CF_PIPELINE_NAME --build-id=$CF_BUILD_ID &
        sleep 5s
		`
	cmd := GetCommandsFactory().CreateWaitRolloutCMD("test", "context")
	if cmd != expectedCMD {
		t.Error("Wait rollout cmd is wrong")
	}
}
