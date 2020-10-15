package builder

import (
	"fmt"
	"net/url"
)

type SyncArgs struct {
	Sync            bool
	WaitHealthy     bool
	WaitForSuspend  bool
	Debug           bool
	AdditionalFlags string
}

type RolloutArgs struct {
	KubernetesContext string
	RolloutName       string
	RolloutNamespace  string
	WaitHealthy       bool
}

type Builder interface {
	Auth(host string, username string, password string, additionalFlags string) error
	Sync(args *SyncArgs, name string, authToken string, host string)
	ExportExternalUrl(host string, name string)
	Rollout(args *RolloutArgs, name string, authToken string, host string)

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

func (b *builder) Auth(host string, username string, password string, additionalFlags string) error {
	domain, err := getHostDomain(host)
	if err != nil {
		return err
	}
	b.lines = append(b.lines, wrapArgoCommandWithAdditionalFlags(fmt.Sprintf("argocd login \"%s\" --insecure --username \"%s\" --password \"%s\"", *domain, username, password), additionalFlags))
	return nil
}

func wrapArgoCommandWithAdditionalFlags(command string, additionalFlags string) string {
	if additionalFlags != "" {
		return command + " " + additionalFlags
	}
	return command
}

func wrapArgoCommandWithToken(command string, authToken string, host string) string {
	if authToken != "" {
		return fmt.Sprintf(command+"  --auth-token %s --server %s --insecure", authToken, host)
	}
	return command
}

func (b *builder) Sync(args *SyncArgs, name string, authToken string, host string) {
	hostDomain, _ := getHostDomain(host)
	var basicCommand string
	if args.Sync {
		basicCommand = fmt.Sprintf("argocd app sync %s", name)
	}
	if args.WaitHealthy {
		basicCommand = fmt.Sprintf("argocd app wait %s", name)
	}
	if args.WaitForSuspend {
		basicCommand = fmt.Sprintf("argocd app wait %s --suspended", name)
	}
	commandWithToken := wrapArgoCommandWithToken(basicCommand, authToken, *hostDomain)
	b.lines = append(b.lines, wrapArgoCommandWithAdditionalFlags(commandWithToken, args.AdditionalFlags))
}

func (b *builder) Rollout(args *RolloutArgs, name string, authToken string, host string) {
	hostDomain, _ := getHostDomain(host)
	b.lines = append(b.lines, "kubectl config get-contexts")
	b.lines = append(b.lines, fmt.Sprintf("kubectl config use-context \"%s\"", args.KubernetesContext))
	b.lines = append(b.lines, fmt.Sprintf("kubectl argo rollouts promote \"%s\" -n \"%s\"", args.RolloutName, args.RolloutNamespace))
	if args.WaitHealthy {
		b.lines = append(b.lines, wrapArgoCommandWithToken(fmt.Sprintf("argocd app wait %s", name), authToken, *hostDomain))
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
