FROM golang:1.14 as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY ./ ./
RUN go build -o ./cf-argo-plugin

FROM alpine
RUN apk --update add curl bash
RUN curl -L https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-amd64 -o /usr/local/bin/kubectl-argo-rollouts
RUN chmod +x /usr/local/bin/kubectl-argo-rollouts

RUN curl -L https://storage.googleapis.com/kubernetes-release/release/v1.17.4/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

RUN curl -sSL https://github.com/argoproj/argo-cd/releases/download/v1.8.5/argocd-linux-amd64 -o /usr/local/bin/argocd
RUN chmod +x /usr/local/bin/argocd

COPY --from=build /app/cf-argo-plugin /usr/local/bin/cf-argo-plugin
ENTRYPOINT /bin/bash