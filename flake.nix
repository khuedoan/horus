{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
      ];
    in
    {
      devShells = supportedSystems (system: {
        default =
          with nixpkgs.legacyPackages.${system};
          mkShell {
            packages = [
              age
              ansible
              ansible-lint
              gnumake
              go
              k3d
              kubectl
              nixfmt-rfc-style
              openssh
              opentofu
              oras
              pre-commit
              shellcheck
              sops
              temporal-cli
              terragrunt
              wireguard-tools
              yamlfmt
              yamllint

              (python3.withPackages (
                p: with p; [
                  kubernetes
                ]
              ))
            ];
          };
      });
    };
}
