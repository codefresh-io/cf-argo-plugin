package builder

import (
	"reflect"
	"testing"
)

func TestBuildTokenFlags(t *testing.T) {
	tokenParams := buildTokenFlags("token", "host", true)
	if tokenParams != " --auth-token token --server host --insecure --prune" {
		t.Errorf("Wrong \"%s\" token command ", tokenParams)
	}

}

func TestBuildTokenFlagsWithoutPrune(t *testing.T) {
	tokenParams := buildTokenFlags("token", "host", false)
	if tokenParams != " --auth-token token --server host --insecure" {
		t.Errorf("Wrong \"%s\" token command ", tokenParams)
	}
}

func TestRolloutWithoutWaitHealthy(t *testing.T) {
	builder := New()
	builder.Rollout(&RolloutArgs{
		KubernetesContext:   "kube-ctx",
		RolloutName:         "app",
		RolloutNamespace:    "default",
		WaitHealthy:         false,
		WaitAdditionalFlags: "",
		Debug:               false,
	}, "test", "token", "host", "context", false)

	expectedLines := []string{
		"#!/bin/bash -e",
		"kubectl config get-contexts",
		"kubectl config use-context \"kube-ctx\"",
		"kubectl argo rollouts promote \"app\" -n \"default\"",
	}

	lines := builder.GetLines()

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Error("Rollout commands is incorrect")
	}

}

func TestRolloutWithWaitHealthy(t *testing.T) {
	builder := New()
	builder.Rollout(&RolloutArgs{
		KubernetesContext:   "kube-ctx",
		RolloutName:         "app",
		RolloutNamespace:    "default",
		WaitHealthy:         true,
		WaitAdditionalFlags: "",
		Debug:               false,
	}, "test", "token", "host", "context", false)

	expectedLines := []string{
		"#!/bin/bash -e",
		"kubectl config get-contexts",
		"kubectl config use-context \"kube-ctx\"",
		"kubectl argo rollouts promote \"app\" -n \"default\"",
		`
		cf-argo-plugin wait-rollout test  --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=context --pipeline-id="$CF_PIPELINE_NAME" --build-id=$CF_BUILD_ID &
        sleep 5s
		`,
		"argocd app wait test   --auth-token token --server  --insecure",
	}

	lines := builder.GetLines()

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Error("Rollout commands is incorrect")
	}

}

func TestSyncWithoutWaitHealthy(t *testing.T) {
	builder := New()
	builder.Sync(&SyncArgs{
		WaitHealthy:         false,
		WaitAdditionalFlags: "",
		Debug:               false,
		Sync:                true,
	}, "test", "token", "host", "context", false)

	expectedLines := []string{
		"#!/bin/bash -e",
		"argocd app sync test  --auth-token token --server  --insecure",
	}

	lines := builder.GetLines()

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Error("Sync commands is incorrect")
	}

}

func TestSyncWithWaitHealthy(t *testing.T) {
	builder := New()
	builder.Sync(&SyncArgs{
		WaitHealthy:         true,
		WaitAdditionalFlags: "",
		Debug:               false,
		Sync:                true,
		Rollback:            true,
	}, "test", "token", "host", "context", false)

	expectedLines := []string{
		"#!/bin/bash -e",
		`
		cf-argo-plugin wait-rollout test  --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=context --pipeline-id="$CF_PIPELINE_NAME" --build-id=$CF_BUILD_ID &
        sleep 5s
		`,
		"argocd app sync test  --auth-token token --server  --insecure",
		`
		{
           set +e
           argocd app wait test   --auth-token token --server  --insecure 2> /codefresh/volume/sync_error.log
        }
        if [[ $? -ne 0 ]]; then
		  ARGO_SYNC_ERROR=$(cat /codefresh/volume/sync_error.log | grep -i fatal)
		  ARGO_SYNC_FAILED=1
		fi
		echo ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"
		cf_export ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"

        wait
        `,
	}

	lines := builder.GetLines()

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Error("Sync commands is incorrect")
	}

}

func TestSyncWithWaitHealthyAndSkip(t *testing.T) {
	builder := New()
	builder.Sync(&SyncArgs{
		WaitHealthy:         true,
		WaitAdditionalFlags: "",
		Debug:               false,
		Sync:                true,
	}, "test", "token", "host", "context", true)

	expectedLines := []string{
		"#!/bin/bash -e",
		`
		cf-argo-plugin wait-rollout test --skip --cf-host=$CF_URL --cf-token=$CF_API_KEY --cf-integration=context --pipeline-id="$CF_PIPELINE_NAME" --build-id=$CF_BUILD_ID &
        sleep 5s
		`,
		"argocd app sync test  --auth-token token --server  --insecure",
		`
		{
           set +e
           argocd app wait test   --auth-token token --server  --insecure 2> /codefresh/volume/sync_error.log
        }
        if [[ $? -ne 0 ]]; then
		  ARGO_SYNC_ERROR=$(cat /codefresh/volume/sync_error.log | grep -i fatal)
		  ARGO_SYNC_FAILED=1
		fi
		echo ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"
		cf_export ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"

        wait
        `,
		`        if [[ -v ARGO_SYNC_FAILED ]]; then
		  exit 1
        fi`,
	}

	lines := builder.GetLines()

	if !reflect.DeepEqual(expectedLines, lines) {
		t.Error("Sync commands is incorrect")
	}

}
