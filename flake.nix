{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      with nixpkgs.legacyPackages.${system};
      {
        devShells.default = mkShell {
          packages = [
            ansible
            ansible-lint
            gnumake
            go
            k3d
            kubectl
            openssh
            opentofu
            pre-commit
            shellcheck
            temporal-cli
            terragrunt
            wireguard-tools
            yamlfmt
            yamllint

            (python3.withPackages (p: with p; [
              kubernetes
            ]))
          ];
        };
      }
    );
}
