{
  description = "Swadloon";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
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
          pname = "swadloon";
          version = self.shortRev or "dirty";
          src = ./.;
          vendorHash = "sha256-ypoSM5rkTb3NF6gDYoBgYWEPkn1kxTPn1CvAGTJ9l3E=";
        };
      in
      {
        packages.default = program;
        packages.swadloon = program;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      }
    );
}
