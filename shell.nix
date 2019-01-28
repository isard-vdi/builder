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

{ pkgs ? import <nixpkgs> {}, goPath ? "go_1_11" }:

with pkgs;
let
  go = lib.getAttrFromPath(lib.SplitString "." goPath) pkgs;
  drv = import ./default.nix { inherit pkgs goPath; };
in

drv.overrideAttrs (attrs: {
  src = null;
  buildInputs = [ govers ] ++ attrs.buildInputs;
  shellHook = ''
    echo "Entering ${attrs.name}"
    set -v

    export GOPATH="$(pwd)/.go"
    export GOCACHE=""
    export GO111MODULE="on"
    go mod download

    set +v
    clear
  '';
})

