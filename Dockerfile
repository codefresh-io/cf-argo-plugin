FROM golang:1.14 as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY ./ ./
RUN go build -o ./cf-argo-plugin

FROM debian:bullseye-slim
RUN apt-get update -y && apt-get install curl bash -y  \
    && curl -L https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-amd64 -o /usr/local/bin/kubectl-argo-rollouts \
    && chmod +x /usr/local/bin/kubectl-argo-rollouts \
    && curl -L https://storage.googleapis.com/kubernetes-release/release/v1.17.4/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    && curl -sSL https://github.com/argoproj/argo-cd/releases/download/v2.4.8/argocd-linux-amd64 -o /usr/local/bin/argocd \
    && chmod +x /usr/local/bin/argocd \
    && apt-get install busybox -y && ln -s /bin/busybox /usr/bin/[[ \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/cf-argo-plugin /usr/local/bin/cf-argo-plugin
ENTRYPOINT /bin/bash