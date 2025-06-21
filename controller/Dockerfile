FROM docker.io/golang:1.24.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/worker ./worker

FROM docker.io/nixos/nix

RUN echo "experimental-features = flakes nix-command" >> /etc/nix/nix.conf

# TODO use native nix develop, currently it's a bit slow
RUN nix-env --install --quiet --attr \
    nixpkgs.kubernetes-helm \
    nixpkgs.opentofu \
    nixpkgs.oras \
    nixpkgs.terragrunt

COPY --from=builder /bin/worker /bin/worker

CMD ["/bin/worker"]
