FROM docker.io/alpine:3.16

RUN apk add --no-cache \
    ansible \
    git \
    make \
    neovim \
    openssh \
    py3-pip \
    python3 \
    terraform=~1.2

RUN apk add --no-cache --repository "http://dl-cdn.alpinelinux.org/alpine/edge/testing" \
    helm \
    kubectl
