FROM golang:1.14 as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY ./ ./
RUN go build -o ./cf-argo-plugin

FROM debian:bullseye-20230502-slim
RUN apt-get update -y && apt-get install wget bash -y  \
    && wget -O /usr/local/bin/kubectl-argo-rollouts https://github.com/argoproj/argo-rollouts/releases/download/v1.5.0/kubectl-argo-rollouts-linux-amd64 \
    && chmod +x /usr/local/bin/kubectl-argo-rollouts \
    && wget -O /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.17.4/bin/linux/amd64/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    && wget -O /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/v2.7.1/argocd-linux-amd64 \
    && chmod +x /usr/local/bin/argocd \
    && apt-get install busybox -y && ln -s /bin/busybox /usr/bin/[[ \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/cf-argo-plugin /usr/local/bin/cf-argo-plugin
ENTRYPOINT /bin/bash