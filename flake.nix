{
  description = "Devshell for crator";

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

        crator = pkgs.buildGoModule {
          pname = "crator";
          version = self.shortRev or "dirty";
          src = ./.;

          vendorHash = "sha256-+D5jLcFWr5djg36xaiHzPFPnZ6XFMPrr+QAj3WA/Yq8=";
        };
      in
      {
        packages.default = crator;
        packages.crator = crator;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      }
    );
}
