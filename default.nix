# Copyright (C) 2019 Néfix Estrada <nefixestrada@gmail.com>
# Author: Néfix Estrada <nefixestrada@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

{ pkgs ? import <nixpkgs> {}, goPath ? "go_1_11", ... }:

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
