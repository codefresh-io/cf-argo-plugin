package builder

import (
	"fmt"
	"net/url"
)

type SyncArgs struct {
	Sync           bool
	WaitHealthy    bool
	WaitForSuspend bool
}

type RolloutArgs struct {
	KubernetesContext string
	RolloutName       string
	RolloutNamespace  string
	WaitHealthy       bool
}

type Builder interface {
	Auth(host string, username string, password string) error
	Sync(args *SyncArgs, name string)
	ExportExternalUrl(host string, name string)
	Rollout(args *RolloutArgs, name string)

	GetLines() []string
}

type builder struct {
	lines []string
}

func New() Builder {
	return &builder{lines: []string{"#!/bin/bash -e"}}
}

func (b *builder) Auth(host string, username string, password string) error {
	domain, err := getHostDomain(host)
	if err != nil {
		return err
	}
	b.lines = append(b.lines, fmt.Sprintf("argocd login \"%s\" --insecure --username \"%s\" --password \"%s\"", *domain, username, password))

	return nil
}

func (b *builder) Sync(args *SyncArgs, name string) {
	if args.Sync {
		b.lines = append(b.lines, fmt.Sprintf("argocd app sync %s", name))
	}
	if args.WaitHealthy {
		b.lines = append(b.lines, fmt.Sprintf("argocd app wait %s", name))
	}
	if args.WaitForSuspend {
		b.lines = append(b.lines, fmt.Sprintf("argocd app wait %s --suspended", name))
	}
}

func (b *builder) Rollout(args *RolloutArgs, name string) {
	b.lines = append(b.lines, "kubectl config get-contexts")
	b.lines = append(b.lines, fmt.Sprintf("kubectl config use-context \"%s\"", args.KubernetesContext))
	b.lines = append(b.lines, fmt.Sprintf("kubectl argo rollouts promote \"%s\" -n \"%s\"", args.RolloutName, args.RolloutNamespace))
	if args.WaitHealthy {
		b.lines = append(b.lines, fmt.Sprintf("argocd app wait %s", name))
	}
}

func (b *builder) GetLines() []string {
	return b.lines
}

func (b *builder) ExportExternalUrl(host string, name string) {
	applicationUrl := fmt.Sprintf("%s/applications/%s", host, name)
	b.lines = append(b.lines, fmt.Sprintf("cf_export runArgoCd_CF_OUTPUT_URL=\"%s\"", applicationUrl))
}

func getHostDomain(host string) (*string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &u.Host, nil
}
