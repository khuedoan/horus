{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
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
            cue
            git
            gnumake
            go
            k9s
            kubectl
            neovim
            openssh
            opentofu
            pre-commit
            shellcheck
            timoni
            wireguard-tools
            yamllint

            (python3.withPackages (p: with p; [
              kubernetes
            ]))
          ];
        };
      }
    );
}
