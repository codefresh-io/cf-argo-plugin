package processResult

import (
	"cf-argo-plugin/pkg/codefresh"
	"testing"
)

type MockCodefresh struct {
}

func (c *MockCodefresh) GetIntegration(name string) (*codefresh.ArgoIntegration, error) {
	return nil, nil
}

func (c *MockCodefresh) StartSyncTask(name string) (*codefresh.TaskResult, error) {
	return nil, nil
}

func (c *MockCodefresh) SendMetadata(metadata *codefresh.ArgoApplicationMetadata) (error, []codefresh.UpdatedActivity) {
	return nil, nil
}

func (c *MockCodefresh) RollbackToStable(name string, payload codefresh.Rollback) (*codefresh.TaskResult, error) {
	return nil, nil
}

func (c *MockCodefresh) GetEnvironments() ([]codefresh.CFEnvironment, error) {
	return nil, nil
}

type MockArgo struct {
}

func (a *MockArgo) GetLatestHistoryId(application string) (int64, error) {
	return 0, nil
}

func TestHandleWaitRollout(t *testing.T) {
	handler := &WaitRolloutHandler{codefresh: &MockCodefresh{}, argo: &MockArgo{}}
	e := handler.processNewHistoryId(123, "test", "123")
	if e == nil {
		t.Error("Should fail with error")
	}
}
