FROM docker.io/golang:1.24.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/worker ./worker

FROM docker.io/nixos/nix

RUN nix-env --install --quiet --attr \
    nixpkgs.terragrunt \
    nixpkgs.opentofu

COPY --from=builder /bin/worker /bin/worker

CMD ["/bin/worker"]
