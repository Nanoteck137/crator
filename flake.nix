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

          vendorHash = "sha256-Hng6WZng+64wCvNq3hBQlLAsexAfU+ifFxeUjVDyofk=";
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
