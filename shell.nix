# https://status.nixos.org (nixos-22.11)
{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/a558f7ac29f5.tar.gz") {} }:

let
  python-packages = pkgs.python3.withPackages (p: with p; [
    kubernetes
  ]);
in
pkgs.mkShell {
  buildInputs = with pkgs; [
    ansible
    ansible-lint
    git
    gnumake
    k9s
    kubectl
    neovim
    oci-cli
    openssh
    pre-commit
    shellcheck
    terraform
    yamllint

    python-packages
  ];
}
