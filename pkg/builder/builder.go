package builder

import (
	"fmt"
	"net/url"
)

type SyncArgs struct {
	Sync                bool
	WaitHealthy         bool
	WaitForSuspend      bool
	Debug               bool
	Prune               bool
	AdditionalFlags     string
	Revision            string
	WaitAdditionalFlags string
	SkipWaitRollout     bool
}

type RolloutArgs struct {
	KubernetesContext   string
	RolloutName         string
	RolloutNamespace    string
	WaitHealthy         bool
	WaitAdditionalFlags string
	Debug               bool
	SkipWaitRollout     bool
}

type Builder interface {
	Auth(host string, username string, password string) error
	Sync(args *SyncArgs, name string, authToken string, host string, context string, skip bool)
	ExportExternalUrl(host string, name string)
	Rollout(args *RolloutArgs, name string, authToken string, host string, context string, skip bool)

	GetLines() []string
	GetExportLines() []string
}

type builder struct {
	lines       []string
	exportLines []string
}

func New() Builder {
	return &builder{lines: []string{"#!/bin/bash -e"}, exportLines: []string{"#!/bin/bash -e"}}
}

func (b *builder) Auth(host string, username string, password string) error {
	domain, err := getHostDomain(host)
	if err != nil {
		return err
	}
	b.lines = append(b.lines, fmt.Sprintf("argocd login \"%s\" --insecure --username \"%s\" --password \"%s\"", *domain, username, password))
	return nil
}

func buildTokenFlags(authToken string, host string, prune bool) string {
	cmd := ""
	if authToken != "" {
		cmd += fmt.Sprintf(" --auth-token %s --server %s --insecure", authToken, host)
	}
	if prune {
		cmd += " --prune"
	}
	return cmd
}

func (b *builder) Sync(args *SyncArgs, name string, authToken string, host string, context string, skip bool) {
	hostDomain, _ := getHostDomain(host)
	tokenFlags := buildTokenFlags(authToken, *hostDomain, args.Prune)
	tokenFlagsWithoutPrune := buildTokenFlags(authToken, *hostDomain, false)
	if args.WaitHealthy || args.WaitForSuspend {
		b.lines = append(b.lines, GetCommandsFactory().CreateWaitRolloutCMD(name, context, skip))
	}

	if args.Sync {
		command := fmt.Sprintf("argocd app sync %s %s", name, tokenFlags)
		if args.Revision != "" {
			command = fmt.Sprintf("%s --revision %s", command, args.Revision)
		}
		b.lines = append(b.lines, command)
	}

	if args.WaitHealthy {
		cmd := fmt.Sprintf(`
		{
           set +e
           argocd app wait %s %s %s 2> /codefresh/volume/sync_error.log
        }
        if [[ $? -ne 0 ]]; then
		  ARGO_SYNC_ERROR=$(cat /codefresh/volume/sync_error.log | grep -i fatal)
		fi
		echo ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"
		cf_export ARGO_SYNC_ERROR="$ARGO_SYNC_ERROR"

        wait
        `, name, args.WaitAdditionalFlags, tokenFlagsWithoutPrune)
		b.lines = append(b.lines, cmd)
	}
	if args.WaitForSuspend {
		b.lines = append(b.lines, fmt.Sprintf("argocd app wait %s %s --suspended", name, tokenFlagsWithoutPrune))
	}
}

func (b *builder) Rollout(args *RolloutArgs, name string, authToken string, host string, context string, skip bool) {
	hostDomain, _ := getHostDomain(host)
	b.lines = append(b.lines, "kubectl config get-contexts")
	b.lines = append(b.lines, fmt.Sprintf("kubectl config use-context \"%s\"", args.KubernetesContext))
	b.lines = append(b.lines, fmt.Sprintf("kubectl argo rollouts promote \"%s\" -n \"%s\"", args.RolloutName, args.RolloutNamespace))

	if args.WaitHealthy {
		if context != "" {
			b.lines = append(b.lines, GetCommandsFactory().CreateWaitRolloutCMD(name, context, skip))
		}
		tokenFlags := buildTokenFlags(authToken, *hostDomain, false)
		b.lines = append(b.lines, fmt.Sprintf("argocd app wait %s %s %s", name, args.WaitAdditionalFlags, tokenFlags))
	}
}

func (b *builder) GetLines() []string {
	return b.lines
}

func (b *builder) GetExportLines() []string {
	return b.exportLines
}

func (b *builder) ExportExternalUrl(host string, name string) {
	applicationUrl := fmt.Sprintf("%s/applications/%s", host, name)
	command := fmt.Sprintf("cf_export runArgoCd_CF_OUTPUT_URL=\"%s\"", applicationUrl)
	b.lines = append(b.lines, command)
	b.exportLines = append(b.exportLines, command)
}

func getHostDomain(host string) (*string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &u.Host, nil
}
