{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  };

  outputs =
    { self, nixpkgs }:
    let
      forAllSystems =
        function:
        nixpkgs.lib.genAttrs [
          "x86_64-linux"
          "aarch64-linux"
        ] (system: function (import nixpkgs { inherit system; }));
    in
    {
      devShells = forAllSystems (pkgs: {
        default =
          with pkgs;
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
            ];
          };
      });
    };
}
