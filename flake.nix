{
  description = "Drup";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url  = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        program = pkgs.buildGoModule {
          pname = "drup";
          version = self.shortRev or "dirty";
          src = ./.;
          vendorHash = "sha256-/EXKoXiejgmx4G/yJDNMd+wKhfflmvtyVnnCdzd8BgI=";
        };
      in
      {
        packages.default = program;
        packages.drup = program;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];
        };
      }
    );
}
