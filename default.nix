{ pkgs ? import <nixpkgs> {}, goPath ? "go_1_11" }:

with pkgs;
let
  go = lib.getAttrFromPath(lib.splitString "." goPath) pkgs;
  buildGoPackage = pkgs.buildGoPackage.override { inherit go; };
in

buildGoPackage {
  name = "isard-builder";
  version = "1.0.0";
  goPackagePath = "github.com/isard-vdi/builder";
  src = lib.cleanSourceWith {
    filter = (path: type:
      ! (builtins.any
      (r: (builtins.match r (builtins.baseNameOf path)) != null)
      [
        ".env"
        ".go"
      ])
    );
    src = lib.cleanSource ./.;
  };
  goDeps = ./deps.nix;
}
